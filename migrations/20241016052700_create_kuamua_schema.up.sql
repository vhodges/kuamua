
BEGIN;

CREATE TABLE IF NOT EXISTS kuamua_patterns
(
    id               bigint GENERATED ALWAYS AS IDENTITY,

    pattern_name     text not null,  -- Callers name for the pattern, returned as the match from Quamina
    pattern          text not null,  -- The pattern to match against

	  group_name       text not null,  -- Patterns are grouped into sets with one Quamina instance per set
	  sub_group_name   text not null,  -- Set to "" if not using

    owner_id         text not null   -- Account or Company id Or Something
);

-- Function to notify listeners that the table has been update

CREATE OR REPLACE FUNCTION notify_updated_kuamua_patterns() RETURNS TRIGGER as
'
  BEGIN
    IF (TG_OP = ''DELETE'') THEN
    PERFORM pg_notify(''kuamua_pattern_updates'', OLD.owner_id || '':'' || OLD.group_name || '':'' || OLD.sub_group_name);
    
    RETURN OLD;
  ELSE
    PERFORM pg_notify(''kuamua_pattern_updates'', NEW.owner_id || '':'' || NEW.group_name || '':'' || NEW.sub_group_name);
    RETURN NEW;
  END IF;
END;
' LANGUAGE plpgsql;

-- Trigger that run the above function
CREATE TRIGGER updated_kuamua_patterns_trigger
AFTER INSERT OR UPDATE OR DELETE ON kuamua_patterns FOR EACH ROW
EXECUTE PROCEDURE notify_updated_kuamua_patterns();


ALTER TABLE kuamua_patterns ADD CONSTRAINT km_owneridlen CHECK (char_length(owner_id) > 0);  -- Can''t be empty, minimum 1

ALTER TABLE kuamua_patterns ADD CONSTRAINT km_namelen CHECK (char_length(pattern_name) > 0); -- Can''t be empty, minimum 1
ALTER TABLE kuamua_patterns ADD CONSTRAINT km_patternlen CHECK (char_length(pattern) > 0);   -- Can''t be empty, minimum 1
ALTER TABLE kuamua_patterns ADD CONSTRAINT km_grouplen CHECK (char_length(group_name) > 0);  -- Can''t be empty, minimum 1
-- Note: SubGroupName CAN be ""

CREATE UNIQUE INDEX IF NOT EXISTS km_unique_name_idx ON kuamua_patterns (pattern_name, group_name, sub_group_name, owner_id);
CREATE INDEX IF NOT EXISTS km_owner_idx ON kuamua_patterns (owner_id);
CREATE INDEX IF NOT EXISTS km_group_idx ON kuamua_patterns (group_name);
CREATE INDEX IF NOT EXISTS km_subgroup_idx ON kuamua_patterns (sub_group_name);

COMMIT;
