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
    UNIQUE KEY (email)
)ENGINE=INNODB;

CREATE TABLE investments (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    InvestmentDate TIMESTAMP,
    StartupName VARCHAR(255) NOT NULL,
    Industry VARCHAR(255) NOT NULL,
    Headquarters VARCHAR(255) NOT NULL,
    BoardRepresentation VARCHAR(255),
    BoardMembers VARCHAR(255),
    CapTable VARCHAR(255),
    InvestmentBackground VARCHAR(255) NOT NULL,
    InvestmentThesis VARCHAR(255) NOT NULL,
    ExitValueAtClosing DECIMAL(10,2),
    FundOwnershipPercentage DECIMAL(10,2),
    InvestorGroupPercentage DECIMAL(10,2),
    ManagementOwnership DECIMAL(10,2),
    InvestmentCommittment DECIMAL(10,2),
    InvestedCapital DECIMAL(10,2),
    RealizedProceeds DECIMAL(10,2),
    ReportedValue DECIMAL(10,2),
    InvestmentMultiple DECIMAL(10,2),
    GrossIRR DECIMAL(10,2),
    UNIQUE KEY (StartupName)
)ENGINE=INNODB;


CREATE TABLE news (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    NewsDate timestamp,
    News VARCHAR(255),
    INDEX news_ind (investment_id),
    FOREIGN KEY (investment_id)
    REFERENCES investments(id)
    ON DELETE CASCADE
)ENGINE=INNODB;


CREATE TABLE investment_structure (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    ReportingDate timestamp,
    Units DECIMAL(10,2),
    TotalInvested DECIMAL(10,2),
    ReportedValue DECIMAL(10,2),
    RealizedProceeds DECIMAL(10,2),
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
    ClosingValue DECIMAL(10,2),
    YearEndValue DECIMAL(10,2),
    	INDEX cs_ind (investment_id),
    	FOREIGN KEY (investment_id)
        REFERENCES investments(id)
        ON DELETE CASCADE
)ENGINE=INNODB;

CREATE TABLE financial_results (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    investment_id bigint(20) unsigned NOT NULL,
    ReportingDate timestamp,
    Revenue DECIMAL(10,2),
    YoYGrowthPercentage1 DECIMAL(10,2),
    LTMEBITDA DECIMAL(10,2),
    YoYGrowthPercentage2 DECIMAL(10,2),
    EBITDAMargin DECIMAL(10,2),
    TotalExitValue DECIMAL(10,2),
    TotalExitValueMultiple DECIMAL(10,2),
    TotalLeverage DECIMAL(10,2),
    TotalLeverageMultiple DECIMAL(10,2),
    	INDEX fr_ind (investment_id),
    	FOREIGN KEY (investment_id)
        REFERENCES investments(id)
        ON DELETE CASCADE
)ENGINE=INNODB;