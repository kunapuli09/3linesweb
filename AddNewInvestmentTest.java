package test.java.com.starpath.smart;

import com.starpath.smart.base.BaseTest;
import org.junit.*;
import org.junit.runner.RunWith;
import org.junit.runners.BlockJUnit4ClassRunner;
import org.openqa.selenium.By;
import org.openqa.selenium.Platform;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeDriverService;
import org.openqa.selenium.remote.DesiredCapabilities;
import org.openqa.selenium.remote.RemoteWebDriver;

import java.io.File;
import java.io.IOException;
import java.net.URL;
import java.util.List;
import java.util.concurrent.TimeUnit;

@RunWith(BlockJUnit4ClassRunner.class)
public class AddNewInvestmentTest extends BaseTest {
    private WebDriver driver;

    private static final By LoginHeadingPath = By.xpath("//*[@id=\"login\"]/div/div[1]/div/h2");
    private static final By LoginEmailFieldPath = By.xpath("//*[@id=\"Email\"]");
    private static final By LoginPasswordFieldPath = By.xpath("//*[@id=\"Password\"]");
    private static final By LoginButtonPath = By.xpath("//*[@id=\"loginButton\"]");
    private static final By AddNewInvestmentButtonPath = By.xpath("//*[@id=\"page-top\"]/section/div/p[2]/a[1]");
    private static final By StartupNameHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[1]/div/div[1]/label");
    private static final By StartupNameFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[1]/div/div[1]/input");
    private static final By TeamHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[1]/div/div[2]/label");
    private static final By TeamFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[1]/div/div[2]/input");
    private static final By WebsiteHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[2]/div/div[1]/label");
    private static final By WebsiteFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[2]/div/div[1]/input");
    private static final By LogoHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[2]/div/div[2]/label");
    private static final By LogoFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[2]/div/div[2]/input");
    private static final By DescriptionFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[3]/div/div/div");
    private static final By IndustryHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[1]/label");
    private static final By IndustryFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[1]/input");
    private static final By InvestmentDateHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/label");
    //InvestmentDate is a date picker
    private static final By InvestmentDateFieldPath = By.xpath("//*[@id=\"InvestmentDate\"]");
    private static final By HeadquartersHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[3]/label");
    private static final By HeadquartersFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[3]/input");
    private static final By CapTableHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[4]/label");
    private static final By CapTableFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[4]/input");
    private static final By BoardRepHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[5]/label");
    private static final By BoardRepFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[5]/input");
    private static final By BoardMembersHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[6]/label");
    private static final By BoardMembersFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[1]/div[6]/input");
    private static final By ExitValueAtClosingHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[1]/label");
    private static final By ExitValueAtClosingFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[1]/input");
    private static final By FundOwnershipPercentageHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[2]/label");
    private static final By FundOwnershipPercentageFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[2]/input");
    private static final By InvestorGroupPercentageHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[3]/label");
    private static final By InvestorGroupPercentageFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[3]/input");
    private static final By ManagementOwnershipPercentageHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[4]/label");
    private static final By ManagementOwnershipPercentageFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[4]/input");
    private static final By InvestmentCommitmentHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[5]/label");
    private static final By InvestmentCommitmentFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[5]/input");
    private static final By InvestedCapitalHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[6]/label");
    private static final By InvestedCapitalFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[2]/div[6]/input");
    private static final By RealizedProceedsHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[1]/label");
    private static final By RealizedProceedsFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[1]/input");
    private static final By ReportedValueHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[2]/label");
    private static final By ReportedValueFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[2]/input");
    private static final By InvestmentMultipleHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[3]/label");
    private static final By InvestmentMultipleFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[3]/input");
    private static final By GrossIRRHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[4]/label");
    private static final By GrossIRRFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[4]/input");
    //Status is a drop-down menu
    private static final By StatusHeadingPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[4]/div[3]/div[5]/label");
    private static final By StatusPendingPath = By.xpath("//*[@id=\"Status\"]/option[1]");
    private static final By StatusCompletePath = By.xpath("//*[@id=\"Status\"]/option[2]");
    private static final By StatusArchivePath = By.xpath("//*[@id=\"Status\"]/option[3]");
    private static final By InvestmentBackgroundFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[5]/div/div/div");
    private static final By InvestmentThesisFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[6]/div/div/div");
    private static final By ValuationMethodologyFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[7]/div/div/div");
    private static final By RiskAssessmentFieldPath = By.xpath("//*[@id=\"newinvestment\"]/div/div[2]/div[8]/div/div/div");
    private static final By LogoutButtonPath = By.xpath("//*[@id=\"navbarCollapse\"]/ul/li[2]/a");

    private WebElement startup, team, website, logo, description, industry, investmentDate, headquarters, capTable;
    private WebElement boardRep, boardMembers, exitValueAtClosing, fundOwnership, investorGroup, managementOwnership;
    private WebElement investmentCommitment, investedCapital, realizedProceeds, reportedValue, investmentMultiple;
    private WebElement grossIRR, statusPending, statusComplete, statusArchive, investmentBackground, investmentThesis;
    private WebElement valuationMethodology, riskAssessment, logoutButton;


    @Before
    public void createAndStartService() {
        System.setProperty("webdriver.chrome.driver", "C:\\Software\\chromedriver_win32\\chromedriver.exe");
        driver = new ChromeDriver();
        driver.get("http://localhost:8888/login");

        //Web Element Declarations
        WebElement loginHeading = driver.findElement(LoginHeadingPath);
        Assert.assertNotNull(loginHeading);
        Assert.assertEquals("LOGIN", loginHeading.getText());

        WebElement loginButton = driver.findElement(LoginButtonPath);
        Assert.assertNotNull(loginButton);

        WebElement email = driver.findElement(LoginEmailFieldPath);
        Assert.assertNotNull(email);

        //Login Elements
        WebElement password = driver.findElement(LoginPasswordFieldPath);
        Assert.assertNotNull(password);

        email.sendKeys("krishna.kunapuli@threelines.us");
        password.sendKeys("Vibhuthi@2018");
        loginButton.click();

        WebElement addNewInvestmentButton = driver.findElement(AddNewInvestmentButtonPath);
        Assert.assertNotNull(addNewInvestmentButton);

        addNewInvestmentButton.click();

        //Set-Up for Add New Investment Page

        WebElement startupHeading = driver.findElement(StartupNameHeadingPath);
        Assert.assertNotNull(startupHeading);
        Assert.assertEquals("Startup Name", startupHeading.getText());

        startup = driver.findElement(StartupNameFieldPath);
        Assert.assertNotNull(startup);

        WebElement teamHeading = driver.findElement(TeamHeadingPath);
        Assert.assertNotNull(teamHeading);
        Assert.assertEquals("Team", teamHeading.getText());

        team = driver.findElement(TeamFieldPath);
        Assert.assertNotNull(team);

        WebElement websiteHeading = driver.findElement(WebsiteHeadingPath);
        Assert.assertNotNull(websiteHeading);
        Assert.assertEquals("Website", websiteHeading.getText());

        website = driver.findElement(WebsiteFieldPath);
        Assert.assertNotNull(website);

        WebElement logoHeading = driver.findElement(LogoHeadingPath);
        Assert.assertNotNull(logoHeading);
        Assert.assertEquals("Logo Name/Path", logoHeading.getText());

        logo = driver.findElement(LogoFieldPath);
        Assert.assertNotNull(logo);

        description = driver.findElement(DescriptionFieldPath);
        Assert.assertNotNull(description);

        WebElement industryHeading = driver.findElement(IndustryHeadingPath);
        Assert.assertNotNull(industryHeading);
        Assert.assertEquals("Industry", industryHeading.getText());

        industry = driver.findElement(IndustryFieldPath);
        Assert.assertNotNull(industry);

        WebElement investmentDateHeading = driver.findElement(InvestmentDateHeadingPath);
        Assert.assertNotNull(investmentDateHeading);
        Assert.assertEquals("Investment Date", investmentDateHeading.getText());

        investmentDate = driver.findElement(InvestmentDateFieldPath);
        Assert.assertNotNull(investmentDate);

        WebElement headquartersHeading = driver.findElement(HeadquartersHeadingPath);
        Assert.assertNotNull(headquartersHeading);
        Assert.assertEquals("Headquarters", headquartersHeading.getText());

        headquarters = driver.findElement(HeadquartersFieldPath);
        Assert.assertNotNull(headquarters);

        WebElement capTableHeading = driver.findElement(CapTableHeadingPath);
        Assert.assertNotNull(capTableHeading);
        Assert.assertEquals("Cap Table", capTableHeading.getText());

        capTable = driver.findElement(CapTableFieldPath);
        Assert.assertNotNull(capTable);

        WebElement boardRepHeading = driver.findElement(BoardRepHeadingPath);
        Assert.assertNotNull(boardRepHeading);
        Assert.assertEquals("Board Representation", boardRepHeading.getText());

        boardRep = driver.findElement(BoardRepFieldPath);
        Assert.assertNotNull(boardRep);

        WebElement boardMembersHeading = driver.findElement(BoardMembersHeadingPath);
        Assert.assertNotNull(boardMembersHeading);
        Assert.assertEquals("Board Members", boardMembersHeading.getText());

        boardMembers = driver.findElement(BoardMembersFieldPath);
        Assert.assertNotNull(boardMembers);

        WebElement exitValueAtClosingHeading = driver.findElement(ExitValueAtClosingHeadingPath);
        Assert.assertNotNull(exitValueAtClosingHeading);
        Assert.assertEquals("Exit Value At Closing", exitValueAtClosingHeading.getText());

        exitValueAtClosing = driver.findElement(ExitValueAtClosingFieldPath);
        Assert.assertNotNull(exitValueAtClosing);

        WebElement fundOwnershipHeading = driver.findElement(FundOwnershipPercentageHeadingPath);
        Assert.assertNotNull(fundOwnershipHeading);
        Assert.assertEquals("Fund Ownership Percentage", fundOwnershipHeading.getText());

        fundOwnership = driver.findElement(FundOwnershipPercentageFieldPath);
        Assert.assertNotNull(fundOwnership);

        WebElement investorGroupHeading = driver.findElement(InvestorGroupPercentageHeadingPath);
        Assert.assertNotNull(investorGroupHeading);
        Assert.assertEquals("Investor Group Percentage", investorGroupHeading.getText());

        investorGroup = driver.findElement(InvestorGroupPercentageFieldPath);
        Assert.assertNotNull(investorGroup);

        WebElement managementOwnershipHeading = driver.findElement(ManagementOwnershipPercentageHeadingPath);
        Assert.assertNotNull(managementOwnershipHeading);
        Assert.assertEquals("Management Ownership Percentage", managementOwnershipHeading.getText());

        managementOwnership = driver.findElement(ManagementOwnershipPercentageFieldPath);
        Assert.assertNotNull(managementOwnership);

        WebElement investmentCommitmentHeading = driver.findElement(InvestmentCommitmentHeadingPath);
        Assert.assertNotNull(investmentCommitmentHeading);
        Assert.assertEquals("Investment Commitment", investmentCommitmentHeading.getText());

        investmentCommitment = driver.findElement(InvestmentCommitmentFieldPath);
        Assert.assertNotNull(investmentCommitment);

        WebElement investedCapitalHeading = driver.findElement(InvestedCapitalHeadingPath);
        Assert.assertNotNull(investedCapitalHeading);
        Assert.assertEquals("Invested Capital", investedCapitalHeading.getText());

        investedCapital = driver.findElement(InvestedCapitalFieldPath);
        Assert.assertNotNull(investedCapital);

        WebElement realizedProceedsHeading = driver.findElement(RealizedProceedsHeadingPath);
        Assert.assertNotNull(realizedProceedsHeading);
        Assert.assertEquals("Realized Proceeds", realizedProceedsHeading.getText());

        realizedProceeds = driver.findElement(RealizedProceedsFieldPath);
        Assert.assertNotNull(realizedProceeds);

        WebElement reportedValueHeading = driver.findElement(ReportedValueHeadingPath);
        Assert.assertNotNull(reportedValueHeading);
        Assert.assertEquals("Reported Value", reportedValueHeading.getText());

        reportedValue = driver.findElement(ReportedValueFieldPath);
        Assert.assertNotNull(reportedValue);

        WebElement investmentMultipleHeading = driver.findElement(InvestmentMultipleHeadingPath);
        Assert.assertNotNull(investmentMultipleHeading);
        Assert.assertEquals("Investment Multiple", investmentMultipleHeading.getText());

        investmentMultiple = driver.findElement(InvestmentMultipleFieldPath);
        Assert.assertNotNull(investmentMultiple);

        WebElement grossIRRHeading = driver.findElement(GrossIRRHeadingPath);
        Assert.assertNotNull(grossIRRHeading);
        Assert.assertEquals("Gross Internal Rate of Return", grossIRRHeading.getText());

        grossIRR = driver.findElement(GrossIRRFieldPath);
        Assert.assertNotNull(grossIRR);

        WebElement statusHeading = driver.findElement(StatusHeadingPath);
        Assert.assertNotNull(statusHeading);
        Assert.assertEquals("Status", statusHeading.getText());

        statusPending = driver.findElement(StatusPendingPath);
        Assert.assertNotNull(statusPending);

        statusComplete = driver.findElement(StatusCompletePath);
        Assert.assertNotNull(statusComplete);

        statusArchive = driver.findElement(StatusArchivePath);
        Assert.assertNotNull(statusArchive);

        investmentBackground = driver.findElement(InvestmentBackgroundFieldPath);
        Assert.assertNotNull(investmentBackground);

        investmentThesis = driver.findElement(InvestmentThesisFieldPath);
        Assert.assertNotNull(investmentThesis);

        valuationMethodology = driver.findElement(ValuationMethodologyFieldPath);
        Assert.assertNotNull(valuationMethodology);

        riskAssessment = driver.findElement(RiskAssessmentFieldPath);
        Assert.assertNotNull(riskAssessment);

        logoutButton = driver.findElement(LogoutButtonPath);
        Assert.assertNotNull(logoutButton);
    }

    @After
    public void createAndStopService() {
        driver.close();
    }


    @Test
    public void testAddNewInvestmentPass() {
        startup.sendKeys("Test Tech");
        team.sendKeys("Name One, Name Two, Name Three");
        website.sendKeys("test-tech.com");
        logo.sendKeys("/img.png");
        description.sendKeys("");
        industry.sendKeys("");

    }

    @Test
    public void testAddNewInvestmentFail() {

    }

    @Test
    public void testNavbar(){

    }

    @Test
    public void testBackToPortfolio(){

    }

    @Test
    public void testFooter(){

    }
}
