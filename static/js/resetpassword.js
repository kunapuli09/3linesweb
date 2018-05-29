$(function() {

  $("#resetForm input,#resetForm textarea").jqBootstrapValidation({
    preventSubmit: true,
    submitError: function($form, event, errors) {
      // additional error messages or events
    },
    submitSuccess: function($form, event) {
      event.preventDefault(); // prevent default submit behaviour
      // get values from FORM
      var email = $("input#ResetEmail").val();
      var firstName = email; // For Success/Failure Message
      // Check for white space in name for Success/Fail message
      if (firstName.indexOf(' ') >= 0) {
        firstName = name.split(' ').slice(0, -1).join(' ');
      }
      $this = $("#resetButton");
      $this.prop("disabled", true); // Disable submit button until AJAX call is complete to prevent duplicate messages
      $.ajax({
        url: "/resetEmail",
        type: "POST",
        data: {
          ResetEmail: email
        },
        cache: false,
        success: function() {
          // Success message
          $('#result').html("<div class='alert alert-success'>");
          $('#result > .alert-success').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#result > .alert-success')
            .append("<strong>Your password reset request has been sent. Pls check your email including spam folder</strong>");
          $('#result > .alert-success')
            .append('</div>');
          //clear all fields
          $('#resetForm').trigger("reset");
        },
        error: function() {
          // Fail message
          $('#result').html("<div class='alert alert-danger'>");
          $('#result > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#result > .alert-danger').append($("<strong>").text("Sorry " + firstName + ", it seems that email is not recognized!"));
          $('#result > .alert-danger').append('</div>');
          //clear all fields
          $('#resetForm').trigger("reset");
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
