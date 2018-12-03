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
import org.openqa.selenium.remote.DesiredCapabilities;
import org.openqa.selenium.remote.RemoteWebDriver;

import java.io.IOException;
import java.net.URL;
import java.util.List;
import java.util.concurrent.TimeUnit;

@RunWith(BlockJUnit4ClassRunner.class)
public class ContactFormTest extends BaseTest {

    private WebDriver driver;
    private WebElement contactFormLink;
    private WebElement name;
    private WebElement email;
    private WebElement phone;
    private WebElement message;
    private WebElement contactFormSubmitButton;

    private static final By NAVBAR_CONTACT_LINK_PATH = By.xpath("//*[@id=\"navbarResponsive\"]/ul/li[5]/a");
    private static final By CONTACT_FORM_HEADING_PATH = By.xpath("//*[@id=\"contact\"]/div[1]/div/div/h2");
    private static final By CONTACT_FORM_NAME_PATH = By.xpath("//*[@id=\"name\"]");
    private static final By CONTACT_FORM_EMAIL_PATH = By.xpath("//*[@id=\"email\"]");
    private static final By CONTACT_FORM_PHONE_PATH = By.xpath("//*[@id=\"phone\"]");
    private static final By CONTACT_FORM_MESSAGE_PATH = By.xpath("//*[@id=\"message\"]");
    private static final By CONTACT_FORM_SUBMIT_BUTTON_PATH = By.xpath("//*[@id=\"sendMessageButton\"]");


    @Before
    public void createAndStartService() {
        System.setProperty("webdriver.chrome.driver", "C:\\Software\\chromedriver_win32\\chromedriver.exe");
        driver = new ChromeDriver();
        driver.get("http://localhost:8888/");

        contactFormLink = driver.findElement(NAVBAR_CONTACT_LINK_PATH);
        Assert.assertNotNull(contactFormLink);

        WebElement contactFormHeading = driver.findElement(CONTACT_FORM_HEADING_PATH);
        Assert.assertNotNull(contactFormHeading);
        Assert.assertEquals("CONTACT US", contactFormHeading.getText());

        name = driver.findElement(CONTACT_FORM_NAME_PATH);
        Assert.assertNotNull(name);

        email = driver.findElement(CONTACT_FORM_EMAIL_PATH);
        Assert.assertNotNull(email);

        phone = driver.findElement(CONTACT_FORM_PHONE_PATH);
        Assert.assertNotNull(phone);

        message = driver.findElement(CONTACT_FORM_MESSAGE_PATH);
        Assert.assertNotNull(message);

        contactFormSubmitButton = driver.findElement(CONTACT_FORM_SUBMIT_BUTTON_PATH);
        Assert.assertNotNull(contactFormSubmitButton);
    }

    @After
    public void createAndStopService() {
        driver.close();
    }

    @Test
    public void testContactForm() {
        contactFormLink.click();

        name.sendKeys("Maddie Hatter");
        email.sendKeys("test"+ Math.round(Math.random() * 10000) +"@test.com");
        phone.sendKeys("3033468288");
        message.sendKeys("Hey, what's up?");

        contactFormSubmitButton.click();
    }
}
