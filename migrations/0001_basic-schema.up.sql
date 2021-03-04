# Run "/usr/local/mysql/bin/mysql 3linesweb -uroot -p < migrations/0001_basic-schema.up.sql"

SET SESSION time_zone = "+0:00";
ALTER DATABASE CHARACTER SET "utf8";



DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS financial_results;
DROP TABLE IF EXISTS capital_structure;
DROP TABLE IF EXISTS investment_structure;
DROP TABLE IF EXISTS news;
DROP TABLE IF EXISTS investments;

CREATE TABLE users (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(12) NOT NULL,
    admin TINYINT(1),
    UNIQUE KEY (email)
)ENGINE=INNODB;

CREATE TABLE investments (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    InvestmentDate TIMESTAMP,
    StartupName VARCHAR(255) NOT NULL,
    FundLegalName VARCHAR(255) NOT NULL,
    Website VARCHAR(255),
    LogoPath VARCHAR(255),
    Description TEXT,
    Team VARCHAR(255),
    RiskAssessment TEXT,
    ValuationMethodology TEXT,
    Industry VARCHAR(255),
    Headquarters VARCHAR(255),
    BoardRepresentation VARCHAR(255),
    BoardMembers VARCHAR(255),
    CapTable VARCHAR(255),
    InvestmentBackground TEXT,
    InvestmentThesis TEXT,
    ExitValueAtClosing DECIMAL(20,2),
    FundOwnershipPercentage DECIMAL(20,2),
    InvestorGroupPercentage DECIMAL(20,2),
    ManagementOwnership DECIMAL(20,2),
    InvestmentCommittment DECIMAL(20,2),
    InvestedCapital DECIMAL(20,2),
    TotalCapitalRaised DECIMAL(20,2),
    RealizedProceeds DECIMAL(20,2),
    ReportedValue DECIMAL(20,2),
    InvestmentMultiple DECIMAL(20,2),
    GrossIRR DECIMAL(20,2),
    Status VARCHAR(255),
    Fund VARCHAR(255),
    UNIQUE KEY (StartupName)
)ENGINE=INNODB;

CREATE TABLE news (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    NewsDate DATE,
    Title VARCHAR(255),
    Status VARCHAR(255),
    News TEXT,
    INDEX news_ind (investment_id),
    FOREIGN KEY (investment_id)
    REFERENCES investments(id)
    ON DELETE CASCADE
)ENGINE=INNODB;

CREATE TABLE assessments (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    ReviewDate timestamp,
    RevenueGrowth TEXT,
    Execution TEXT,
    Leadership TEXT,
    RevenueBreakEvenPlan TEXT,
    KeyGrowthEnablers TEXT,
    PlaybookAdoption VARCHAR(255),
    Status VARCHAR(255),
    StartupName VARCHAR(255) NOT NULL,
    MarketMultiple DECIMAL(20,2),
    YearThreeForecastedRevenue DECIMAL(20,2),
    ThreelinesValueAtExit DECIMAL(20,2),
    YearThreeExitMultiple DECIMAL(20,2),
    INDEX review_ind (investment_id),
        FOREIGN KEY (investment_id)
        REFERENCES investments(id)
        ON DELETE CASCADE
)ENGINE=INNODB;

CREATE TABLE docs (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    UploadDate DATE,
    DocPath VARCHAR(255),
    DocName VARCHAR(255),
    Hash VARCHAR(255), 
    INDEX docs_ind (investment_id),
    FOREIGN KEY (investment_id)
    REFERENCES investments(id)
    ON DELETE CASCADE
)ENGINE=INNODB;

CREATE TABLE userdocs (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id bigint(20) unsigned NOT NULL,
    UploadDate DATE,
    DocPath VARCHAR(255),
    DocName VARCHAR(255),
    Hash VARCHAR(255), 
    INDEX docs_ind (user_id),
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
)ENGINE=INNODB;


CREATE TABLE investment_structure (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    ReportingDate timestamp,
    Units DECIMAL(20,2),
    TotalInvested DECIMAL(20,2),
    ReportedValue DECIMAL(20,2),
    RealizedProceeds DECIMAL(20,2),
    Structure VARCHAR(255),
    	INDEX is_ind (investment_id),
    	FOREIGN KEY (investment_id)
        REFERENCES investments(id)
        ON DELETE CASCADE
)ENGINE=INNODB;

CREATE TABLE capital_structure (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    ReportingDate timestamp,
    ClosingValue DECIMAL(20,2),
    YearEndValue DECIMAL(20,2),
    Capitalization VARCHAR(255),
    	INDEX cs_ind (investment_id),
    	FOREIGN KEY (investment_id)
        REFERENCES investments(id)
        ON DELETE CASCADE
)ENGINE=INNODB;

CREATE TABLE financial_results (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    ReportingDate timestamp,
    Revenue DECIMAL(20,2),
    YoYGrowthPercentage1 DECIMAL(20,2),
    LTMEBITDA DECIMAL(20,2),
    YoYGrowthPercentage2 DECIMAL(20,2),
    EBITDAMargin DECIMAL(20,2),
    TotalExitValue DECIMAL(20,2),
    TotalExitValueMultiple DECIMAL(20,2),
    TotalLeverage DECIMAL(20,2),
    TotalLeverageMultiple DECIMAL(20,2),
    Assessment VARCHAR(255),
    	INDEX fr_ind (investment_id),
    	FOREIGN KEY (investment_id)
        REFERENCES investments(id)
        ON DELETE CASCADE
)ENGINE=INNODB;

CREATE TABLE applications (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ApplicationDate timestamp,
    FirstName VARCHAR(255) NOT NULL,
    LastName VARCHAR(255) NOT NULL,
    Email VARCHAR(255) NOT NULL,
    CompanyName VARCHAR(255) NOT NULL,
    Website VARCHAR(255) NOT NULL,
    Phone VARCHAR(12) NOT NULL,
    Title VARCHAR(255) NOT NULL,
    Referrer VARCHAR(255) NOT NULL,
    Industries VARCHAR(255),
    Locations VARCHAR(255),
    Revenue VARCHAR(255),
    CapitalRaised DECIMAL(20,2),
    Comments TEXT,
    ElevatorPitch TEXT(255),
    UNIQUE KEY (email),
    UNIQUE KEY (Website)
)ENGINE=INNODB;

CREATE TABLE screeningnotes (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    application_id bigint(20) unsigned NOT NULL,
    ScreeningDate timestamp,
    ScreenerEmail VARCHAR(255) NOT NULL,
    Need VARCHAR(255) NOT NULL,
    Status VARCHAR(255) NOT NULL,
    TeamRisk TINYINT NOT NULL DEFAULT '5',
    BarrierToEntry TINYINT NOT NULL DEFAULT '5',
    TechRisk TINYINT NOT NULL DEFAULT '5',
    CompetitionRisk TINYINT NOT NULL DEFAULT '5',
    PoliticalRisk TINYINT NOT NULL DEFAULT '5',
    SupplierRisk TINYINT NOT NULL DEFAULT '5',
    ExecutionRisk TINYINT NOT NULL DEFAULT '5',
    MarketRisk TINYINT NOT NULL DEFAULT '5',
    ScalingRisk TINYINT NOT NULL DEFAULT '5',
    ExitRisk TINYINT NOT NULL DEFAULT '5',
    Comments TEXT,
    INDEX scr_ind (application_id),
        FOREIGN KEY (application_id)
        REFERENCES applications(id)
        ON DELETE CASCADE
)ENGINE=INNODB;


CREATE TABLE notifications (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    news_id bigint(20) unsigned NOT NULL,
    NotificationDate timestamp DEFAULT CURRENT_TIMESTAMP,
    StartupName VARCHAR(255) NOT NULL,
    Status VARCHAR(255) NOT NULL,
    Email VARCHAR(255) NOT NULL,
    Title VARCHAR(255) NOT NULL,
    NewsDate timestamp DEFAULT CURRENT_TIMESTAMP
)ENGINE=INNODB;

CREATE TABLE contributions (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id bigint(20) unsigned NOT NULL,
    FundLegalName VARCHAR(255) NOT NULL,
    InvestorLegalName VARCHAR(255) NOT NULL,
    InvestorAddress VARCHAR(255) NOT NULL,
    InvestorType VARCHAR(255) NOT NULL,
    GroupContact VARCHAR(255),
    InvestmentGroupName VARCHAR(255),
    CommitmentDate timestamp,
    OwnershipPercentage DECIMAL(20,2),
    InvestmentAmount DECIMAL(20,2),
    Comments VARCHAR(255),
    Status VARCHAR(255),
    INDEX user_ind (user_id),
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
)ENGINE=INNODB;


CREATE TABLE `offices` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `FirstName` varchar(255) NOT NULL,
  `LastName` varchar(255) NOT NULL,
  `EmailAddress` varchar(255) NOT NULL,
  `FirmType` varchar(255) DEFAULT NULL,
  `FirmName` varchar(255) DEFAULT NULL,
  `Address1` varchar(255) DEFAULT NULL,
  `Address2` varchar(255) DEFAULT NULL,
  `City` varchar(255) DEFAULT NULL,
  `State` varchar(255) DEFAULT NULL,
  `Zip Code` varchar(255) DEFAULT NULL,
  `Phone` varchar(255) DEFAULT NULL,
  `Fax` varchar(255) DEFAULT NULL,
  `Website` varchar(255) DEFAULT NULL,
  `Description` text,
  `ContactName` varchar(255) NOT NULL,
  `ContactTitle` text,
  `ContactPhone` varchar(255) DEFAULT NULL,
  `ContactLocation` varchar(255) DEFAULT NULL,
  `ContactSpecialty` varchar(255) DEFAULT NULL,
  `Stages` text,
  `PortfolioFirms` text,
  `RecentFundings` text,
  `FirmFocus` text,
  `LastUpdatedBy` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=34514 DEFAULT CHARSET=latin1

CREATE TABLE executives (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    ApplicationDate timestamp,
    Name VARCHAR(255) NOT NULL,
    Email VARCHAR(255) NOT NULL,
    SocialMediaHandle TEXT NOT NULL,
    UNIQUE KEY (email)
)ENGINE=INNODB;


alter table assessments ADD StartupName VARCHAR(255) NOT NULL;
alter table assessments ADD MarketMultiple DECIMAL(20,2);
alter table assessments ADD YearThreeForecastedRevenue DECIMAL(20,2);
alter table assessments ADD ThreelinesValueAtExit DECIMAL(20,2);
alter table assessments ADD YearThreeExitMultiple DECIMAL(20,2);

alter table applications ADD  Referrer VARCHAR(255);
alter table applications ADD  ElevatorPitch VARCHAR(255);
alter table applications ADD  Revenue VARCHAR(255);
alter table applications ADD  Status VARCHAR(255);
ALTER TABLE applications ADD INDEX  Status_Index (Status);
ALTER TABLE applications ADD LastUpdatedBy VARCHAR(255);
ALTER TABLE applications ADD LastUpdatedTime timestamp;

ALTER TABLE applications ADD INDEX  CompanyName_Index (CompanyName);
ALTER TABLE applications ADD INDEX  Locations_Index (Locations);
ALTER TABLE investments ADD FundLegalName VARCHAR(255);
ALTER TABLE investments ADD TotalCapitalRaised DECIMAL(20,2)
SET SQL_SAFE_UPDATES = 0;
UPDATE  investments SET FundLegalName="3Lines 2016 Discretionary Fund, LLC";
#Roles update for new blogreaders
ALTER TABLE 3linesweb.users ADD Roles VARCHAR(255); 
ALTER TABLE 3linesweb.users DROP COLUMN Admin;
UPDATE 3linesweb.users set Roles='Admin,Dsc,Investor,BlogReader' where admin=1;

UPDATE 3linesweb.users set Roles='Investor,BlogReader' where admin=0;

UPDATE 3linesweb.users set Roles='Dsc,Investor,BlogReader' where email in ("fundone@3lines.vc",
        "roy.rajiv@gmail.com",
        "arun.taman@gmail.com",
        "sgosala99@gmail.com",
        "dsc@3lines.vc")

DROP TABLE 3linesweb.proposaldocs;
#alter table investments ADD  Status VARCHAR(255);
#update investments set Status = "COMPLETE";

#alter table news ADD  Status VARCHAR(255);
#update news set Status = "PUBLISHED";
#alter table news ADD  Title VARCHAR(255);
#update news set Title = "Title 1";