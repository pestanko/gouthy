-- This file should undo anything in `up.sql`-- DROP TRIGGER IF EXISTS SetUpdated_Users ON Users;
-- DROP TRIGGER IF EXISTS SetUpdated_UserSecrets ON UserSecrets;
-- DROP INDEX IF EXISTS Users_Username_Idx;

-- Users

DROP TRIGGER IF EXISTS SetUpdated_Users ON Users;
DROP TRIGGER IF EXISTS SetUpdated_UserSecrets ON UserSecrets;
DROP INDEX IF EXISTS Users_Username_Idx;
DROP TABLE IF EXISTS Users CASCADE ;
DROP TABLE IF EXISTS UserSecrets CASCADE ;


-- Applications
DROP TRIGGER IF EXISTS SetUpdated_Applications ON Applications;
DROP TRIGGER IF EXISTS SetUpdated_ApplicationSecrets ON ApplicationSecrets;

DROP INDEX IF EXISTS Applications_codename_Idx;
DROP INDEX IF EXISTS Applications_client_id_Idx;

DROP TABLE IF EXISTS Applicationsn CASCADE ;
DROP TABLE IF EXISTS ApplicationSecrets CASCADE ;

