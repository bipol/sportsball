--create table
CREATE SCHEMA `sportsball` ;

--create manager table
CREATE TABLE `sportsball`.`manager` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `full_name` VARCHAR(45) NOT NULL,
  `team_id` INT NULL,
  PRIMARY KEY (`id`)
  FOREIGN KEY (team_id)
      REFERENCES team(id)
      ON DELETE CASCADE
  );

--create team table
CREATE TABLE `sportsball`.`team` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `manager` BIGINT NULL,
  PRIMARY KEY (`id`)
  FOREIGN KEY (manager)
      REFERENCES manager(id)
      ON DELETE CASCADE
  );

--create player table
CREATE TABLE `sportsball`.`player` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `full_name` VARCHAR(45) NOT NULL,
  `team_id` BIGINT NULL,
  `field_position` INT NOT NULL,
  PRIMARY KEY (`id`)
  FOREIGN KEY (team_id)
      REFERENCES team(id)
      ON DELETE CASCADE
  );

--create player transition table
CREATE TABLE `sportsball`.`team_player` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `team` BIGINT NOT NULL,
  `player` BIGINT NOT NULL,
  PRIMARY KEY (`id`)
  FOREIGN KEY (team)
      REFERENCES team(id)
      ON DELETE CASCADE
  FOREIGN KEY (player)
      REFERENCES player(id)
      ON DELETE CASCADE
  );

ALTER TABLE `sportsball`.`manager`
ADD COLUMN `access_key` VARCHAR(45) NOT NULL AFTER `team_id`,
ADD UNIQUE INDEX `access_key_UNIQUE` (`access_key` ASC);
