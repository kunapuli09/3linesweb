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
public class SignupPageTest  extends BaseTest{
    private WebDriver driver;
    private WebElement email;
    private WebElement phone;
    private WebElement password;
    private WebElement passwordAgain;
    private WebElement signupButton;

    //private static final By SIGNUP_HEADING_PATH = By.xpath("//*[@id=\"signup\"]/div/div[1]/div/h2");
    private static final By SIGNUP_EMAIL_TEXT_FIELD_PATH = By.xpath("//*[@id=\"signup\"]/div/div[2]/div/form/input[1]");
    private static final By SIGNUP_PHONE_TEXT_FIELD_PATH = By.xpath("//*[@id=\"Phone\"]");
    private static final By SIGNUP_PASSWORD_TEXT_FIELD_PATH = By.xpath("//*[@id=\"signup\"]/div/div[2]/div/form/input[3]");
    private static final By SIGNUP_PASSWORD_AGAIN_TEXT_FIELD_PATH = By.xpath("//*[@id=\"signup\"]/div/div[2]/div/form/input[4]");
    private static final By SIGNUP_SIGNUP_BUTTON_PATH = By.xpath("//*[@id=\"signup\"]/div/div[2]/div/form/button");
    private static final By LOGOUT_BUTTON_PATH = By.xpath("//*[@id=\"navbarCollapse\"]/ul/li[2]/a");

    @Before
    public void createAndStartService() {
        System.setProperty("webdriver.chrome.driver","C:\\Software\\chromedriver_win32\\chromedriver.exe");
        driver = new ChromeDriver();
        driver.get("http://localhost:8888/signup");

        /*WebElement signupHeadingPath = driver.findElement(SIGNUP_HEADING_PATH);
        Assert.assertNotNull(signupHeadingPath);
        Assert.assertEquals("SIGNUP", signupHeadingPath.getText());*/

        email = driver.findElement(SIGNUP_EMAIL_TEXT_FIELD_PATH);
        Assert.assertNotNull(email);

        phone = driver.findElement(SIGNUP_PHONE_TEXT_FIELD_PATH);
        Assert.assertNotNull(phone);

        password = driver.findElement(SIGNUP_PASSWORD_TEXT_FIELD_PATH);
        Assert.assertNotNull(password);

        passwordAgain = driver.findElement(SIGNUP_PASSWORD_AGAIN_TEXT_FIELD_PATH);
        Assert.assertNotNull(passwordAgain);

        signupButton = driver.findElement(SIGNUP_SIGNUP_BUTTON_PATH);
        Assert.assertNotNull(signupButton);
    }

    @After
    public void createAndStopService() {
        driver.close();
    }

    @Test
    public void testSignupLogout() throws InterruptedException {
        email.sendKeys("tensixthshatterhats@gmail.com");
        phone.sendKeys("0000000000");
        password.sendKeys("tensixths@2018");
        passwordAgain.sendKeys("tensixths@2018");
        signupButton.click();

        WebElement logoutButton = driver.findElement(LOGOUT_BUTTON_PATH);
        Assert.assertNotNull(logoutButton);

        logoutButton.click();
        sleep(5000);
    }
}
