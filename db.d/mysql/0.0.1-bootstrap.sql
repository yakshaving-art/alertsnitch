DROP PROCEDURE IF EXISTS bootstrap;

DELIMITER //
CREATE PROCEDURE bootstrap()
BEGIN
  SET @exists := (SELECT 1 FROM information_schema.tables I WHERE I.table_name = "Model" AND I.table_schema = database());
  IF @exists IS NULL THEN

    CREATE TABLE `Model` (
      `ID` enum('1') NOT NULL,
      `version` VARCHAR(20) NOT NULL,
      PRIMARY KEY (`ID`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

    INSERT INTO `Model` (`version`) VALUES ("0.0.1");

  ELSE
    SIGNAL SQLSTATE '42000' SET MESSAGE_TEXT='Model Table Exists, quitting...';
  END IF;
END;
//
DELIMITER ;

-- Execute the procedure
CALL bootstrap();

-- Drop the procedure
DROP PROCEDURE bootstrap;

-- Create the rest of the tables
CREATE TABLE `AlertGroup` (
	`ID` INT NOT NULL AUTO_INCREMENT,
	`time` TIMESTAMP NOT NULL,
	`receiver` VARCHAR(100) NOT NULL,
	`status` VARCHAR(50) NOT NULL,
	`externalURL` TEXT NOT NULL,
	`groupKey` VARCHAR(255) NOT NULL,
	KEY `idx_time` (`time`) USING BTREE,
    KEY `idx_status_ts` (`status`, `time`) USING BTREE,
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `GroupLabel` (
	`ID` INT NOT NULL AUTO_INCREMENT,
    `AlertGroupID` INT NOT NULL,
    `GroupLabel` VARCHAR(100) NOT NULL,
    `Value` VARCHAR(1000) NOT NULL,
    FOREIGN KEY (AlertGroupID) REFERENCES AlertGroup (ID) ON DELETE CASCADE,
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `CommonLabel` (
	`ID` INT NOT NULL AUTO_INCREMENT,
    `AlertGroupID` INT NOT NULL,
    `Label` VARCHAR(100) NOT NULL,
    `Value` VARCHAR(1000) NOT NULL,
    FOREIGN KEY (AlertGroupID) REFERENCES AlertGroup (ID) ON DELETE CASCADE,
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `CommonAnnotation` (
	`ID` INT NOT NULL AUTO_INCREMENT,
    `AlertGroupID` INT NOT NULL,
    `Annotation` VARCHAR(100) NOT NULL,
    `Value` VARCHAR(1000) NOT NULL,
    FOREIGN KEY (AlertGroupID) REFERENCES AlertGroup (ID) ON DELETE CASCADE,
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `Alert` (
	`ID` INT NOT NULL AUTO_INCREMENT,
    `alertGroupID` INT NOT NULL,
	`status` VARCHAR(50) NOT NULL,
    `startsAt` DATETIME NOT NULL,
    `endsAt` DATETIME DEFAULT NULL,
	`generatorURL` TEXT NOT NULL,
    FOREIGN KEY (alertGroupID) REFERENCES AlertGroup (ID) ON DELETE CASCADE,
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `AlertLabel` (
	`ID` INT NOT NULL AUTO_INCREMENT,
    `AlertID` INT NOT NULL,
    `Label` VARCHAR(100) NOT NULL,
    `Value` VARCHAR(1000) NOT NULL,
    FOREIGN KEY (AlertID) REFERENCES Alert (ID) ON DELETE CASCADE,
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `AlertAnnotation` (
	`ID` INT NOT NULL AUTO_INCREMENT,
    `AlertID` INT NOT NULL,
    `Annotation` VARCHAR(100) NOT NULL,
    `Value` VARCHAR(1000) NOT NULL,
    FOREIGN KEY (AlertID) REFERENCES Alert (ID) ON DELETE CASCADE,
	PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
