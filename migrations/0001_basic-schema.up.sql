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
    RealizedProceeds DECIMAL(20,2),
    ReportedValue DECIMAL(20,2),
    InvestmentMultiple DECIMAL(20,2),
    GrossIRR DECIMAL(20,2),
    Status VARCHAR(255),
    UNIQUE KEY (StartupName)
)ENGINE=INNODB;


CREATE TABLE news (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    NewsDate DATE,
    News TEXT,
    INDEX news_ind (investment_id),
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
    Industries VARCHAR(255),
    Locations VARCHAR(255),
    CapitalRaised DECIMAL(20,2),
    Comments VARCHAR(255),
    UNIQUE KEY (email),
    UNIQUE KEY (Website)
)ENGINE=INNODB;