--
-- PostgreSQL database dump
--

-- Dumped from database version 11.4 (Debian 11.4-1.pgdg80+1)
-- Dumped by pg_dump version 11.4 (Debian 11.4-1.pgdg80+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pajbot2; Type: SCHEMA; Schema: -; Owner: -
--

SET default_with_oids = false;
SET default_tablespace = '';

--
-- Name: banphrase; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE banphrase (
    id bigint NOT NULL,
    group_id bigint,
    enabled boolean DEFAULT true,
    description text,
    phrase text NOT NULL,
    length bigint DEFAULT '60'::bigint,
    warning_id bigint,
    case_sensitive boolean,
    type bigint DEFAULT '0'::bigint,
    sub_immunity boolean DEFAULT false,
    remove_accents boolean DEFAULT false
);


--
-- Name: TABLE banphrase; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON TABLE banphrase IS 'Store banned phrases';


--
-- Name: COLUMN banphrase.enabled; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.enabled IS 'NULL = Inherit from group';


--
-- Name: COLUMN banphrase.description; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.description IS 'Optional description of the banphrase, i.e. racism or banned emote';


--
-- Name: COLUMN banphrase.phrase; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.phrase IS 'The banned phrase itself. This can be a regular expression, it all depends on the "operator" of the banphrase';


--
-- Name: COLUMN banphrase.length; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.length IS 'NULL = Inherit from group, 0 = permaban, >0 = timeout for X seconds';


--
-- Name: COLUMN banphrase.warning_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.warning_id IS 'NULL = Inherit from group, anything else is an ID to a warning "scale"';


--
-- Name: COLUMN banphrase.case_sensitive; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.case_sensitive IS 'NULL = Inherit from group';


--
-- Name: COLUMN banphrase.type; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.type IS 'NULL = Inherit from group, 0 = contains, more IDs can be found in the go code lol xd';


--
-- Name: COLUMN banphrase.sub_immunity; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.sub_immunity IS 'NULL = Inherit from group';


--
-- Name: COLUMN banphrase.remove_accents; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase.remove_accents IS 'NULL = Inherit from group';


--
-- Name: banphrase_group; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE banphrase_group (
    id bigint NOT NULL,
    enabled boolean DEFAULT true NOT NULL,
    name character varying(64) NOT NULL,
    description text,
    length bigint DEFAULT '60'::bigint NOT NULL,
    warning_id bigint,
    case_sensitive boolean DEFAULT false NOT NULL,
    type bigint DEFAULT '0'::bigint NOT NULL,
    sub_immunity boolean DEFAULT false NOT NULL,
    remove_accents boolean DEFAULT false NOT NULL
);


--
-- Name: TABLE banphrase_group; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON TABLE banphrase_group IS 'Store banphrase groups. this will make it easier to manage multiple banphrases at the same time';


--
-- Name: COLUMN banphrase_group.description; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase_group.description IS 'Optional description of the banphrase group, i.e. racism or banned emote';


--
-- Name: COLUMN banphrase_group.length; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase_group.length IS '0 = permaban, >0 = timeout for X seconds';


--
-- Name: COLUMN banphrase_group.warning_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase_group.warning_id IS 'ID to a warning "scale"';


--
-- Name: COLUMN banphrase_group.type; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN banphrase_group.type IS '0 = contains, more IDs can be found in the go code lol xd';


--
-- Name: banphrase_group_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE banphrase_group_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: banphrase_group_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE banphrase_group_id_seq OWNED BY banphrase_group.id;


--
-- Name: banphrase_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE banphrase_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: banphrase_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE banphrase_id_seq OWNED BY banphrase.id;


--
-- Name: bot; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE bot (
    id bigint NOT NULL,
    twitch_userid character varying(64) NOT NULL,
    twitch_username character varying(64) NOT NULL,
    twitch_access_token character varying(64),
    twitch_refresh_token character varying(64),
    twitch_access_token_expiry timestamp with time zone
);


--
-- Name: TABLE bot; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON TABLE bot IS 'Store available bot accounts, requires an access token with chat_login scope';


--
-- Name: COLUMN bot.twitch_access_token; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN bot.twitch_access_token IS 'Bot level access-token';


--
-- Name: COLUMN bot.twitch_refresh_token; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN bot.twitch_refresh_token IS 'Bot level refresh-token';


--
-- Name: bot_channel; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE bot_channel (
    id bigint NOT NULL,
    bot_id bigint NOT NULL,
    twitch_channel_id character varying(64) NOT NULL
);


--
-- Name: COLUMN bot_channel.twitch_channel_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN bot_channel.twitch_channel_id IS 'i.e. 11148817';


--
-- Name: bot_channel_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE bot_channel_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bot_channel_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE bot_channel_id_seq OWNED BY bot_channel.id;


--
-- Name: bot_channel_module; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE bot_channel_module (
    id bigint NOT NULL,
    bot_channel_id bigint NOT NULL,
    module_id character varying(128) NOT NULL,
    enabled boolean,
    settings bytea
);


--
-- Name: COLUMN bot_channel_module.module_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN bot_channel_module.module_id IS 'i.e. nuke';


--
-- Name: COLUMN bot_channel_module.enabled; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN bot_channel_module.enabled IS 'if null, it uses the modules default enabled value';


--
-- Name: COLUMN bot_channel_module.settings; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN bot_channel_module.settings IS 'json blob with settings';


--
-- Name: bot_channel_module_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE bot_channel_module_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bot_channel_module_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE bot_channel_module_id_seq OWNED BY bot_channel_module.id;


--
-- Name: bot_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE bot_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bot_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE bot_id_seq OWNED BY bot.id;



--
-- Name: moderation_action; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE moderation_action (
    id bigint NOT NULL,
    channel_id character varying(64) NOT NULL,
    user_id character varying(64) NOT NULL,
    target_id character varying(64) NOT NULL,
    action smallint NOT NULL,
    duration bigint,
    reason text,
    "timestamp" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    context text
);


--
-- Name: COLUMN moderation_action.channel_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN moderation_action.channel_id IS 'Twitch Channel owners user ID';


--
-- Name: COLUMN moderation_action.user_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN moderation_action.user_id IS 'Source user ID';


--
-- Name: COLUMN moderation_action.target_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN moderation_action.target_id IS 'Target user ID (the user who has banned/unbanned/timed out)';


--
-- Name: COLUMN moderation_action.action; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN moderation_action.action IS 'Action in int format, enums declared outside of SQL';


--
-- Name: COLUMN moderation_action.duration; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN moderation_action.duration IS 'Duration of action (only used for timeouts atm)';


--
-- Name: COLUMN moderation_action.reason; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN moderation_action.reason IS 'Reason for ban. Auto filled in from twich chat, but can be modified in web gui';


--
-- Name: COLUMN moderation_action."timestamp"; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN moderation_action."timestamp" IS 'Timestamp of when the timeout occured';


--
-- Name: moderation_action_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE moderation_action_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: moderation_action_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE moderation_action_id_seq OWNED BY moderation_action.id;


--
-- Name: report; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE report (
    id bigint NOT NULL,
    channel_id character varying(64) NOT NULL,
    channel_name character varying(64) NOT NULL,
    channel_type character varying(64) NOT NULL,
    reporter_id character varying(64) NOT NULL,
    reporter_name character varying(64) NOT NULL,
    target_id character varying(64) NOT NULL,
    target_name character varying(64) NOT NULL,
    reason text,
    logs text,
    "time" timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: report_history; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE report_history (
    id bigint NOT NULL,
    channel_id character varying(64) NOT NULL,
    channel_name character varying(64) NOT NULL,
    channel_type character varying(64) NOT NULL,
    reporter_id character varying(64) NOT NULL,
    reporter_name character varying(64) NOT NULL,
    target_id character varying(64) NOT NULL,
    target_name character varying(64) NOT NULL,
    reason text,
    logs text,
    "time" timestamp with time zone,
    handler_id character varying(64) NOT NULL,
    handler_name character varying(64) NOT NULL,
    action smallint NOT NULL,
    action_duration bigint DEFAULT '0'::bigint NOT NULL,
    time_handled timestamp with time zone
);


--
-- Name: COLUMN report_history.channel_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.channel_id IS 'twitch ID of channel user was reported in';


--
-- Name: COLUMN report_history.channel_name; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.channel_name IS 'twitch username of channel the user was reported in';


--
-- Name: COLUMN report_history.reporter_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.reporter_id IS 'twitch user ID of reporter';


--
-- Name: COLUMN report_history.reporter_name; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.reporter_name IS 'twitch user name of reporter';


--
-- Name: COLUMN report_history.target_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.target_id IS 'twitch user ID of person being reported';


--
-- Name: COLUMN report_history.target_name; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.target_name IS 'twitch user name of person being reported';


--
-- Name: COLUMN report_history."time"; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history."time" IS 'time report was added';


--
-- Name: COLUMN report_history.handler_id; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.handler_id IS 'twitch user ID of person who handled the report';


--
-- Name: COLUMN report_history.handler_name; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.handler_name IS 'twitch user name of person who handled the report';


--
-- Name: COLUMN report_history.action; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.action IS 'number constant for what action was taken for the report. 1 = ban, 2 = timeout, 3 = dismiss';


--
-- Name: COLUMN report_history.action_duration; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON COLUMN report_history.action_duration IS 'number of seconds for the action. only relevant for timeouts';


--
-- Name: report_history_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE report_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: report_history_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE report_history_id_seq OWNED BY report_history.id;


--
-- Name: report_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE report_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: report_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE report_id_seq OWNED BY report.id;


--
-- Name: twitch_user_channel_permission; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE twitch_user_channel_permission (
    twitch_user_id character varying(64) NOT NULL,
    channel_id character varying(64) NOT NULL,
    permissions bigint NOT NULL
);


--
-- Name: twitch_user_global_permission; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE twitch_user_global_permission (
    twitch_user_id character varying(64) NOT NULL,
    permissions bigint NOT NULL
);


--
-- Name: user; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE "user" (
    id bigint NOT NULL,
    twitch_username character varying(64) NOT NULL,
    twitch_userid character varying(64) NOT NULL
);


--
-- Name: user_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE user_id_seq OWNED BY "user".id;


--
-- Name: user_session; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE user_session (
    id character varying(64) NOT NULL,
    user_id bigint NOT NULL,
    expiry_date timestamp with time zone
);


--
-- Name: warning_scale; Type: TABLE; Schema: pajbot2; Owner: -
--

CREATE TABLE warning_scale (
    id bigint NOT NULL
);


--
-- Name: TABLE warning_scale; Type: COMMENT; Schema: pajbot2; Owner: -
--

COMMENT ON TABLE warning_scale IS 'Store data about warning scales';


--
-- Name: warning_scale_id_seq; Type: SEQUENCE; Schema: pajbot2; Owner: -
--

CREATE SEQUENCE warning_scale_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: warning_scale_id_seq; Type: SEQUENCE OWNED BY; Schema: pajbot2; Owner: -
--

ALTER SEQUENCE warning_scale_id_seq OWNED BY warning_scale.id;



--
-- Name: banphrase id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY banphrase ALTER COLUMN id SET DEFAULT nextval('banphrase_id_seq'::regclass);


--
-- Name: banphrase_group id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY banphrase_group ALTER COLUMN id SET DEFAULT nextval('banphrase_group_id_seq'::regclass);


--
-- Name: bot id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY bot ALTER COLUMN id SET DEFAULT nextval('bot_id_seq'::regclass);


--
-- Name: bot_channel id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY bot_channel ALTER COLUMN id SET DEFAULT nextval('bot_channel_id_seq'::regclass);


--
-- Name: bot_channel_module id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY bot_channel_module ALTER COLUMN id SET DEFAULT nextval('bot_channel_module_id_seq'::regclass);


--
-- Name: moderation_action id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY moderation_action ALTER COLUMN id SET DEFAULT nextval('moderation_action_id_seq'::regclass);


--
-- Name: report id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY report ALTER COLUMN id SET DEFAULT nextval('report_id_seq'::regclass);


--
-- Name: report_history id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY report_history ALTER COLUMN id SET DEFAULT nextval('report_history_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY "user" ALTER COLUMN id SET DEFAULT nextval('user_id_seq'::regclass);


--
-- Name: warning_scale id; Type: DEFAULT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY warning_scale ALTER COLUMN id SET DEFAULT nextval('warning_scale_id_seq'::regclass);


--
-- Name: banphrase idx_16414_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY banphrase
    ADD CONSTRAINT idx_16414_primary PRIMARY KEY (id);


--
-- Name: banphrase_group idx_16428_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY banphrase_group
    ADD CONSTRAINT idx_16428_primary PRIMARY KEY (id);


--
-- Name: bot idx_16443_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY bot
    ADD CONSTRAINT idx_16443_primary PRIMARY KEY (id);


--
-- Name: bot_channel idx_16449_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY bot_channel
    ADD CONSTRAINT idx_16449_primary PRIMARY KEY (id);


--
-- Name: bot_channel_module idx_16455_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY bot_channel_module
    ADD CONSTRAINT idx_16455_primary PRIMARY KEY (id);


--
-- Name: moderation_action idx_16464_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY moderation_action
    ADD CONSTRAINT idx_16464_primary PRIMARY KEY (id);


--
-- Name: report idx_16474_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY report
    ADD CONSTRAINT idx_16474_primary PRIMARY KEY (id);


--
-- Name: report_history idx_16484_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY report_history
    ADD CONSTRAINT idx_16484_primary PRIMARY KEY (id);


--
-- Name: user idx_16494_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY "user"
    ADD CONSTRAINT idx_16494_primary PRIMARY KEY (id);


--
-- Name: user_session idx_16498_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY user_session
    ADD CONSTRAINT idx_16498_primary PRIMARY KEY (id);


--
-- Name: warning_scale idx_16503_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY warning_scale
    ADD CONSTRAINT idx_16503_primary PRIMARY KEY (id);


--
-- Name: twitch_user_global_permission idx_16653_primary; Type: CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY twitch_user_global_permission
    ADD CONSTRAINT idx_16653_primary PRIMARY KEY (twitch_user_id);


--
-- Name: idx_16414_group_id; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE INDEX idx_16414_group_id ON banphrase USING btree (group_id);


--
-- Name: idx_16414_warning_id; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE INDEX idx_16414_warning_id ON banphrase USING btree (warning_id);


--
-- Name: idx_16428_group_name; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE UNIQUE INDEX idx_16428_group_name ON banphrase_group USING btree (name);


--
-- Name: idx_16428_warning_id; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE INDEX idx_16428_warning_id ON banphrase_group USING btree (warning_id);


--
-- Name: idx_16443_itwitchuid; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE UNIQUE INDEX idx_16443_itwitchuid ON bot USING btree (twitch_userid);


--
-- Name: idx_16449_bot_channel; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE UNIQUE INDEX idx_16449_bot_channel ON bot_channel USING btree (bot_id, twitch_channel_id);


--
-- Name: idx_16455_bot_channel_module; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE UNIQUE INDEX idx_16455_bot_channel_module ON bot_channel_module USING btree (bot_channel_id, module_id);


--
-- Name: idx_16464_channeltargetaction_index; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE INDEX idx_16464_channeltargetaction_index ON moderation_action USING btree (channel_id, target_id, action);


--
-- Name: idx_16464_channelusertarget_index; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE INDEX idx_16464_channelusertarget_index ON moderation_action USING btree (channel_id, user_id, target_id);


--
-- Name: idx_16494_ui_twitch_userid; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE UNIQUE INDEX idx_16494_ui_twitch_userid ON "user" USING btree (twitch_userid);


--
-- Name: idx_16498_user_id; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE INDEX idx_16498_user_id ON user_session USING btree (user_id);


--
-- Name: idx_16650_user_channel_permission; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE UNIQUE INDEX idx_16650_user_channel_permission ON twitch_user_channel_permission USING btree (twitch_user_id, channel_id);


--
-- Name: idx_16653_user_permission; Type: INDEX; Schema: pajbot2; Owner: -
--

CREATE UNIQUE INDEX idx_16653_user_permission ON twitch_user_global_permission USING btree (twitch_user_id, permissions);


--
-- Name: banphrase banphrase_ibfk_1; Type: FK CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY banphrase
    ADD CONSTRAINT banphrase_ibfk_1 FOREIGN KEY (warning_id) REFERENCES warning_scale(id) ON UPDATE RESTRICT ON DELETE SET NULL;


--
-- Name: banphrase banphrase_ibfk_2; Type: FK CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY banphrase
    ADD CONSTRAINT banphrase_ibfk_2 FOREIGN KEY (group_id) REFERENCES banphrase_group(id) ON UPDATE RESTRICT ON DELETE SET NULL;


--
-- Name: banphrase_group banphrasegroup_ibfk_1; Type: FK CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY banphrase_group
    ADD CONSTRAINT banphrasegroup_ibfk_1 FOREIGN KEY (warning_id) REFERENCES warning_scale(id) ON UPDATE RESTRICT ON DELETE SET NULL;


--
-- Name: bot_channel botchannel_ibfk_1; Type: FK CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY bot_channel
    ADD CONSTRAINT botchannel_ibfk_1 FOREIGN KEY (bot_id) REFERENCES bot(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: bot_channel_module botchannelmodule_ibfk_1; Type: FK CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY bot_channel_module
    ADD CONSTRAINT botchannelmodule_ibfk_1 FOREIGN KEY (bot_channel_id) REFERENCES bot_channel(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- Name: user_session usersession_ibfk_1; Type: FK CONSTRAINT; Schema: pajbot2; Owner: -
--

ALTER TABLE ONLY user_session
    ADD CONSTRAINT usersession_ibfk_1 FOREIGN KEY (user_id) REFERENCES "user"(id) ON UPDATE RESTRICT ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

