$(function() {

  $("#loginForm input,#loginForm textarea").jqBootstrapValidation({
    preventSubmit: true,
    submitError: function($form, event, errors) {
      // additional error messages or events
    },
    submitSuccess: function($form, event) {
      event.preventDefault(); // prevent default submit behaviour
      // get values from FORM
      var email = $("input#Email").val();
      var password = $("input#Password").val();
      
      $this = $("#loginButton");
      $this.prop("disabled", true); // Disable submit button until AJAX call is complete to prevent duplicate messages
      $.ajax({
        url: "/login",
        type: "POST",
        data: {
          Email: email,
          Password: password
        },
        cache: false,
        error: function() {
          // Fail message
          $('#success').html("<div class='alert alert-danger'>");
          $('#success > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#success > .alert-danger').append($("<strong>").text("Password didn't match or email not registered. Pls contact system administrator."));
          $('#success > .alert-danger').append('</div>');
          //clear all fields
          $('#loginForm').trigger("reset");
        },
        success: function() {
          window.location.href = '/portfolio'
        },
        complete: function() {
          setTimeout(function() {
            $this.prop("disabled", false); // Re-enable submit button when AJAX call is complete
          }, 1000);
        }
      });
    },
    filter: function() {
      return $(this).is(":visible");
    },
  });

  $("a[data-toggle=\"tab\"]").click(function(e) {
    e.preventDefault();
    $(this).tab("show");
  });
});
/*When clicking on Full hide fail/success boxes */
$('#name').focus(function() {
  $('#success').html('');
});