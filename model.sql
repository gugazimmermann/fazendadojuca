-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema fazendadojuca
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema fazendadojuca
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `fazendadojuca` DEFAULT CHARACTER SET utf8 ;
USE `fazendadojuca` ;

-- -----------------------------------------------------
-- Table `fazendadojuca`.`animal`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `fazendadojuca`.`animal` (
  `ID` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `number` VARCHAR(255) NOT NULL,
  `registry` VARCHAR(255) NOT NULL,
  `origin` VARCHAR(255) NOT NULL,
  `father` INT NOT NULL,
  `mother` INT NOT NULL,
  `insemination` TINYINT NOT NULL,
  `birth` DATE NOT NULL,
  `death` DATE NOT NULL,
  `gender_id` INT NOT NULL,
  `breed_id` INT NOT NULL,
  `purity_level_id` INT NOT NULL,
  PRIMARY KEY (`ID`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `fazendadojuca`.`gender`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `fazendadojuca`.`gender` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `fazendadojuca`.`breed`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `fazendadojuca`.`breed` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `fazendadojuca`.`purity_level`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `fazendadojuca`.`purity_level` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `level` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

-- -----------------------------------------------------
-- Data for table `fazendadojuca`.`gender`
-- -----------------------------------------------------
START TRANSACTION;
USE `fazendadojuca`;
INSERT INTO `fazendadojuca`.`gender` (`id`, `name`) VALUES (1, 'Macho');
INSERT INTO `fazendadojuca`.`gender` (`id`, `name`) VALUES (2, 'Fêmea');

COMMIT;


-- -----------------------------------------------------
-- Data for table `fazendadojuca`.`breed`
-- -----------------------------------------------------
START TRANSACTION;
USE `fazendadojuca`;
INSERT INTO `fazendadojuca`.`breed` (`id`, `name`) VALUES (1, 'Desconhecida');
INSERT INTO `fazendadojuca`.`breed` (`id`, `name`) VALUES (2, 'Aberdeen Angus');
INSERT INTO `fazendadojuca`.`breed` (`id`, `name`) VALUES (3, 'American Angus');
INSERT INTO `fazendadojuca`.`breed` (`id`, `name`) VALUES (4, 'Black Angus');
INSERT INTO `fazendadojuca`.`breed` (`id`, `name`) VALUES (5, 'Red Angus');

COMMIT;


-- -----------------------------------------------------
-- Data for table `fazendadojuca`.`purity_level`
-- -----------------------------------------------------
START TRANSACTION;
USE `fazendadojuca`;
INSERT INTO `fazendadojuca`.`purity_level` (`id`, `level`) VALUES (1, '1');
INSERT INTO `fazendadojuca`.`purity_level` (`id`, `level`) VALUES (2, '1/2');
INSERT INTO `fazendadojuca`.`purity_level` (`id`, `level`) VALUES (3, '3/4');
INSERT INTO `fazendadojuca`.`purity_level` (`id`, `level`) VALUES (4, '7/8');
INSERT INTO `fazendadojuca`.`purity_level` (`id`, `level`) VALUES (5, '15/16');
INSERT INTO `fazendadojuca`.`purity_level` (`id`, `level`) VALUES (6, '31/32');

COMMIT;

