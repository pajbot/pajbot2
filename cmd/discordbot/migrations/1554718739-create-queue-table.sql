CREATE TABLE discord_queue (
	id SERIAL PRIMARY KEY,
	action TEXT NOT NULL,
	timepoint TIMESTAMP with time zone NOT NULL
);

comment on column discord_queue.action is 'json blob describing the action';
comment on column discord_queue.timepoint is 'when to perform the action';

CREATE INDEX timepoint_idx ON discord_queue (timepoint);
