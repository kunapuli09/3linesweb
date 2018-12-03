package test.java.com.starpath.smart;

import com.starpath.smart.base.BaseTest;
import org.junit.*;
import org.junit.runner.RunWith;
import org.junit.runners.BlockJUnit4ClassRunner;
import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.chrome.ChromeDriver;


import static java.lang.Thread.sleep;

@RunWith(BlockJUnit4ClassRunner.class)
public class ApplicationFormTest extends BaseTest {

    private WebDriver driver;

    private WebElement firstName;
    private WebElement lastName;
    private WebElement email;
    private WebElement companyName;
    private WebElement website;
    private WebElement phone;
    private WebElement title;
    private WebElement artificialIntelligenceCheckbox;
    private WebElement unitedStatesRadioButton;
    private WebElement capitalRaised;
    private WebElement fundingNeeds;
    private WebElement submitAppButton;

    private static final By APP_HEADING_PATH = By.xpath("//*[@id=\"contact\"]/div/div[1]/div/h2");
    private static final By FIRST_NAME_TEXT_FIELD_PATH = By.xpath("//*[@id=\"FirstName\"]");
    private static final By LAST_NAME_TEXT_FIELD_PATH = By.xpath("//*[@id=\"LastName\"]");
    private static final By EMAIL_TEXT_FIELD_PATH = By.xpath("//*[@id=\"Email\"]");
    private static final By COMPANY_NAME_TEXT_FIELD_PATH = By.xpath("//*[@id=\"CompanyName\"]");
    private static final By WEBSITE_TEXT_FIELD_PATH = By.xpath("//*[@id=\"Website\"]");
    private static final By PHONE_TEXT_FIELD_PATH = By.xpath("//*[@id=\"Phone\"]");
    private static final By TITLE_TEXT_FIELD_PATH = By.xpath("//*[@id=\"Title\"]");

    private static final By INDUSTRY_HEADING_PATH = By.xpath("//*[@id=\"Industries\"]/div[1]/h6");
    private static final String artificialIntelligenceID = "check1";
    private static final String saasID = "check2";
    private static final String healthcareID = "check3";
    private static final String cyberSecurityID = "check4";
    private static final String internetOfThingsID = "check5";
    private static final String insuranceTechnologyID = "check6";
    private static final String agricultureID = "check7";

    private static final By LOCATIONS_HEADING_PATH = By.xpath("//*[@id=\"Locations\"]/div[1]/h6");
    private static final String unitedStatesID = "inlineRadio1";
    private static final String indiaID = "inlineRadio2";

    private static final By CAPITAL_RAISED_TEXT_FIELD_PATH = By.xpath("//*[@id=\"CapitalRaised\"]");
    private static final By FUNDING_NEEDS_TEXT_FIELD_PATH = By.xpath("//*[@id=\"Comments\"]");
    private static final By SUBMIT_APPLICATION_BUTTON_PATH = By.xpath("//*[@id=\"sendApplicationButton\"]");
    private static final By SUCCESS_MESSAGE_PATH = By.xpath("//*[@id=\"successApplication\"]/div/strong");

    @Before
    public void createAndStartService(){
        System.setProperty("webdriver.chrome.driver", "C:\\Software\\chromedriver_win32\\chromedriver.exe");
        driver = new ChromeDriver();
        driver.get("http://localhost:8888/appl");

        WebElement appHeading = driver.findElement(APP_HEADING_PATH);
        Assert.assertEquals("APPLICATION FOR SEED FUNDING", appHeading.getText());

        firstName = driver.findElement(FIRST_NAME_TEXT_FIELD_PATH);
        Assert.assertNotNull(firstName);

        lastName = driver.findElement(LAST_NAME_TEXT_FIELD_PATH);
        Assert.assertNotNull(lastName);

        email = driver.findElement(EMAIL_TEXT_FIELD_PATH);
        Assert.assertNotNull(email);

        companyName = driver.findElement(COMPANY_NAME_TEXT_FIELD_PATH);
        Assert.assertNotNull(companyName);

        website = driver.findElement(WEBSITE_TEXT_FIELD_PATH);
        Assert.assertNotNull(website);

        phone = driver.findElement(PHONE_TEXT_FIELD_PATH);
        Assert.assertNotNull(phone);

        title = driver.findElement(TITLE_TEXT_FIELD_PATH);
        Assert.assertNotNull(title);

        WebElement industryHeading = driver.findElement(INDUSTRY_HEADING_PATH);
        Assert.assertEquals("Please check the industries most closely related to your company:", industryHeading.getText());

        artificialIntelligenceCheckbox = driver.findElement(By.id(artificialIntelligenceID));
        Assert.assertNotNull(artificialIntelligenceCheckbox);

        WebElement saasCheckbox = driver.findElement(By.id(saasID));
        Assert.assertNotNull(saasCheckbox);

        WebElement healthcareCheckbox = driver.findElement(By.id(healthcareID));
        Assert.assertNotNull(healthcareCheckbox);

        WebElement cyberSecurityCheckbox = driver.findElement(By.id(cyberSecurityID));
        Assert.assertNotNull(cyberSecurityCheckbox);

        WebElement internetOfThingsCheckbox = driver.findElement(By.id(internetOfThingsID));
        Assert.assertNotNull(internetOfThingsCheckbox);

        WebElement insuranceTechnologyCheckbox = driver.findElement(By.id(insuranceTechnologyID));
        Assert.assertNotNull(insuranceTechnologyCheckbox);

        WebElement agricultureCheckbox = driver.findElement(By.id(agricultureID));
        Assert.assertNotNull(agricultureCheckbox);

        WebElement locationsHeading = driver.findElement(LOCATIONS_HEADING_PATH);
        Assert.assertEquals("Please choose one of our locations:", locationsHeading.getText());

        unitedStatesRadioButton = driver.findElement(By.id(unitedStatesID));
        Assert.assertNotNull(unitedStatesRadioButton);

        WebElement indiaRadioButton = driver.findElement(By.id(indiaID));
        Assert.assertNotNull(indiaRadioButton);

        capitalRaised = driver.findElement(CAPITAL_RAISED_TEXT_FIELD_PATH);
        Assert.assertNotNull(capitalRaised);

        fundingNeeds = driver.findElement(FUNDING_NEEDS_TEXT_FIELD_PATH);
        Assert.assertNotNull(fundingNeeds);

        submitAppButton = driver.findElement(SUBMIT_APPLICATION_BUTTON_PATH);
        Assert.assertNotNull(submitAppButton);
    }

    @After
    public void createAndStopService() {
        driver.close();
    }

    @Test
    public void testApplicationForm() throws InterruptedException {
        firstName.sendKeys("Test");
        lastName.sendKeys("One");
        email.sendKeys("test"+ Math.round(Math.random() * 10000) +"@test.com");
        companyName.sendKeys("TestOne");
        website.sendKeys("https://test"+ Math.round(Math.random() * 10000) +".com");
        phone.sendKeys("1111111111");
        title.sendKeys("CEO");

        artificialIntelligenceCheckbox.getAttribute("checked");
        if(artificialIntelligenceCheckbox.equals("false"))
            artificialIntelligenceCheckbox.click();

        unitedStatesRadioButton.getAttribute("checked");
        if(unitedStatesRadioButton.equals("false"))
            unitedStatesRadioButton.click();

        capitalRaised.sendKeys("70000.00");
        System.out.println(capitalRaised.getAttribute("value"));
        fundingNeeds.sendKeys("Gimme the money.");
        submitAppButton.click();

        WebElement successMessage = driver.findElement(SUCCESS_MESSAGE_PATH);
        Assert.assertNotNull(successMessage);

        sleep(1000);
    }

    @Test
    public void testEmptyFieldsError() throws InterruptedException {
        //Comment a line, run, and check. Process repeated for all input lines.
        //Should return false for empty lines.
        firstName.sendKeys("Test");
        lastName.sendKeys("One");
        email.sendKeys("test"+ Math.round(Math.random() * 10000) +"@test.com");
        companyName.sendKeys("TestOne");
        website.sendKeys("https://test"+ Math.round(Math.random() * 10000) +".com");
        phone.sendKeys("1111111111");
        title.sendKeys("CEO");

        artificialIntelligenceCheckbox.getAttribute("checked");
        if(artificialIntelligenceCheckbox.equals("false"))
            artificialIntelligenceCheckbox.click();

        unitedStatesRadioButton.getAttribute("checked");
        if(unitedStatesRadioButton.equals("false"))
            unitedStatesRadioButton.click();

        capitalRaised.sendKeys("70000.00");
        System.out.println(capitalRaised.getAttribute("value"));
        fundingNeeds.sendKeys("Gimme the money.");
        submitAppButton.click();

        WebElement successMessage = driver.findElement(SUCCESS_MESSAGE_PATH);
        Assert.assertNotNull(successMessage);


        sleep(1000);
    }
}
