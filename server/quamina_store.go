package server

import (
	"fmt"
	"os"
    "context"
	//"io"
    "log"
	//"net/http"
	"sync"

	"github.com/vhodges/kuamua/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/hashicorp/golang-lru/v2"
	"quamina.net/go/quamina"

)

type QuaminaStore struct {
	cacheKeyVersions sync.Map

	pool *pgxpool.Pool // Database connection pool
	database *database.Queries
	cache *lru.Cache[string, *quamina.Quamina]
} 

func NewQuaminaStore(pool *pgxpool.Pool, lruCapacity int) (*QuaminaStore, error) {
	
	db := database.New(pool)

	cache, cerr := lru.New[string, *quamina.Quamina](lruCapacity)
	if cerr != nil {
		return nil, cerr
	}

	store := &QuaminaStore{pool: pool, database: db, cache: cache}
	
	// Start the listener
	go func() {
        err := store.cacheBuster()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: Start cacheBuster", err)
			return
		}
    }()

	return store, nil
}

func (store *QuaminaStore) GetQuamina(ownerid string, groupname string, subgroupname string) (*quamina.Quamina, error) {

	cachekey := fmt.Sprintf("%s:%s:%s", ownerid, groupname, subgroupname)
	versionedkey := fmt.Sprintf("%s@%d", cachekey, 
	   store.getCacheKeyVersion(cachekey))

	//fmt.Printf("versionedKey: %s\n", versionedkey)

	var quamina *quamina.Quamina
	var err error
	var ok bool

	quamina, ok = store.cache.Get(versionedkey)

	if !ok {

		//fmt.Printf("***** Cache Miss...\n")
		quamina, err = store.loadQuamina(ownerid, groupname, subgroupname)
		if err != nil {
			return nil, err
		}

		store.cache.Add(versionedkey, quamina)
		return quamina, nil // No need to copy, freshly made
	} 

	//fmt.Printf("     Cache HIT...\n")
	return quamina.Copy(), nil
}

func (store *QuaminaStore) loadQuamina(ownerid string, groupname string, subgroupname string) (*quamina.Quamina, error) {
	quamina, qerr := quamina.New()
    if qerr != nil {
        return nil, qerr
    }

	query := database.ListOwnerGroupPatternsParams{
		OwnerID: ownerid,
		GroupName: groupname,
		SubGroupName: subgroupname,
	}
	
	patterns, err := store.database.ListOwnerGroupPatterns(context.Background(), query)
    if err != nil {
        return nil, err
    }

	for _, pattern := range patterns {
        // Ignore the error, there should be none  TODO Add validation 
        // that the pattern is valid in the Crud create/update routes
        _ = quamina.AddPattern(pattern.PatternName, pattern.Pattern)
	}

    return quamina, nil
}

func (store *QuaminaStore) getCacheKeyVersion(cachekey string) int {
	v, _ := store.cacheKeyVersions.LoadOrStore(cachekey, 0)
	return v.(int)
}

func (store *QuaminaStore) incrementCacheKeyVersion(cachekey string) {
	v, _ := store.cacheKeyVersions.LoadOrStore(cachekey, 0)
	current := v.(int)
	current = current + 1
	store.cacheKeyVersions.Store(cachekey, current)
}

// run in a go routine
func (store *QuaminaStore) cacheBuster() (error) {
	
	conn, err := store.pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "listen kuamua_pattern_updates")
	if err != nil {
		return err
	}

	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error waiting for notification:", err)
		}

		log.Println("Pattern Change, update cache version - PID:", notification.PID, "Channel:", notification.Channel, "Payload:", notification.Payload)

		store.incrementCacheKeyVersion(notification.Payload)
	}

	return nil
}
