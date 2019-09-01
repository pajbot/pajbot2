CREATE TYPE command_action_type AS ENUM('text_response', 'module_action');

CREATE TYPE module_action_id AS ENUM('modules.nuke.nuke', 'modules.message_height_limit.heighttest');
COMMENT ON TYPE module_action_id IS 'Available raw actions that can be bound as the action of module commands';

CREATE TABLE command (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  bot_channel_id INT NOT NULL REFERENCES bot_channel(id),
  cost INT NOT NULL DEFAULT 0,
  cd_user INTERVAL NOT NULL DEFAULT interval '15 seconds',
  cd_all INTERVAL NOT NULL DEFAULT interval '5 seconds',
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  action_type command_action_type NOT NULL, -- this is the descriminator for the polymorphism (type = ENUM('text_response', 'module_action'))
  response TEXT, -- for text_response
  action_id module_action_id -- for module_action
  -- could add CONSTRAINT to check response not null when action_type = 'text_response' for example
);
COMMENT ON TABLE command IS 'Available command on a given bot, in a given channel';

CREATE TABLE command_trigger (
  command_id INT REFERENCES command(id),
  trigger TEXT NOT NULL UNIQUE, -- notice the extra unique constraint
  PRIMARY KEY (command_id, trigger)
);
COMMENT ON TABLE command_trigger IS 'Aliases/Triggers for a given command';

CREATE TYPE permission AS ENUM(
  'twitch_moderator',
  'create_command',
  'edit_command',
  'delete_command',
  'add_banphrase' -- etc.
);
COMMENT ON TYPE permission IS 'Available permissions that can be granted, revoked and required';

CREATE TABLE command_execute_permission_requirement (
  command_id INT REFERENCES command(id),
  permission_id permission NOT NULL,
  PRIMARY KEY (command_id, permission_id)
);
COMMENT ON TABLE command_execute_permission_requirement IS 'Permissions required to execute a command';

