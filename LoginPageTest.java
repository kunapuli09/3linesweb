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
import java.util.concurrent.TimeUnit;

import static java.lang.Thread.*;

@RunWith(BlockJUnit4ClassRunner.class)
public class LoginPageTest extends BaseTest {

    private WebDriver driver;
    private WebElement loginButton;
    private WebElement email;
    private WebElement password;
    //private WebElement forgotPasswordLink;

    private static final By LOGIN_HEADING_PATH = By.xpath("//*[@id=\"login\"]/div/div[1]/div/h2");
    private static final By LOGIN_EMAIL_TEXT_FIELD_PATH = By.xpath("//*[@id=\"Email\"]");
    private static final By LOGIN_PASSWORD_TEXT_FIELD_PATH = By.xpath("//*[@id=\"Password\"]");
    private static final By LOGIN_LOGIN_BUTTON_PATH = By.xpath("//*[@id=\"loginButton\"]");
    private static final By LOGOUT_BUTTON_PATH = By.xpath("//*[@id=\"navbarCollapse\"]/ul/li[2]/a");
    //private static final By FORGOT_PASSWORD_LINK_PATH = By.xpath("//*[@id=\"formLogin\"]/a");


    @Before
    public void createAndStartService(){
        System.setProperty("webdriver.chrome.driver", "C:\\Software\\chromedriver_win32\\chromedriver.exe");
        driver = new ChromeDriver();
        driver.get("http://localhost:8888/login");

        WebElement loginHeadingPath = driver.findElement(LOGIN_HEADING_PATH);
        Assert.assertNotNull(loginHeadingPath);
        Assert.assertEquals("LOGIN", loginHeadingPath.getText());

        loginButton = driver.findElement(LOGIN_LOGIN_BUTTON_PATH);
        Assert.assertNotNull(loginButton);

        email = driver.findElement(LOGIN_EMAIL_TEXT_FIELD_PATH);
        Assert.assertNotNull(email);

        password = driver.findElement(LOGIN_PASSWORD_TEXT_FIELD_PATH);
        Assert.assertNotNull(password);

        /*forgotPasswordLink = driver.findElement(FORGOT_PASSWORD_LINK_PATH);
        Assert.assertNotNull(forgotPasswordLink);*/
    }

    @After
    public void createAndStopService() {
        driver.close();
    }

    @Test
    public void testLoginLogout() throws InterruptedException {

        email.sendKeys("tellme@tellme.com");
        password.sendKeys("tellme@2018");
        loginButton.click();

        WebElement logoutButton = driver.findElement(LOGOUT_BUTTON_PATH);
        Assert.assertNotNull(logoutButton);

        logoutButton.click();

        sleep(5000);
    }

/*    @Test
    public void testForgotPassword() throws InterruptedException {
        forgotPasswordLink.click();

        WebElement forgotPasswordHeadingPath = driver.findElement("//*[@id="pwdModal"]/div/div/div[1]/h1");
        Assert.assertNotNull(forgotPasswordHeadingPath);
        Assert.assertEquals("Forgot Password?", forgotPasswordHeadingPath.getText());

        WebElement forgotPasswordInstructions = driver.findElement("//*[@id="pwdModal"]/div/div/div[2]/div/div/div/div/p");
        Assert.assertNotNull(forgotPasswordInstructions);
        Assert.assertEquals("If you have forgotten your password you can reset it here.", forgotPasswordInstructions.getText());

        WebElement forgotPasswordEmail = driver.findElement(//*[@id="ResetEmail"]);
        Assert.assertNotNull(forgotPasswordEmail);

        forgotPasswordEmail.sendKeys("tensixthshatterhats@gmail.com");

        sleep(5000);
    }*/
}