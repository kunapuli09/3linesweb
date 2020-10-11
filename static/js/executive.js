var onloadCallback3 = function() {
    grecaptcha.render('g-recaptcha3', {
        'sitekey': '6LcF3OQUAAAAAGmMrHmVIWUp4qxjL8wdLnGR6k-w'
    });
  };
$(function() {
  //Check if required fields are filled
    function checkifreqfld() {
        var isFormFilled = true;
        $("#executiveForm").find(".form-control:visible").each(function() {
            var value = $.trim($(this).val());
            if ($(this).prop('required')) {
                if (value.length < 1) {
                    //$(this).closest(".form-group").addClass("field-error");
                    isFormFilled = false;
                } else {
                    //$(this).closest(".form-group").removeClass("field-error");
                }
            } else {
                //$(this).closest(".form-group").removeClass("field-error");
            }
        });
        return isFormFilled;
    }
  $("#executiveForm input,#executiveForm textarea").jqBootstrapValidation({
    preventSubmit: true,
    submitError: function($form, event, errors) {
      // additional error messages or events
    },
    submitSuccess: function($form, event) {
      if (checkifreqfld()) {
                event.preventDefault();
      }
      //google captcha response
      var rcres = grecaptcha.getResponse();
      // get values from FORM
      var Name = $("input#Name").val();
      var Email = $("input#Email").val();
      var SocialMediaHandle = $("input#SocialMediaHandle").val();
      $this = $("#sendExecutiveButton");
      $this.prop("disabled", true); // Disable submit button until AJAX call is complete to prevent duplicate messages
      $.ajax({
        url: "/addExecutive",
        type: "POST",
        data: {
          Name: Name,
          Email: Email,
          SocialMediaHandle: SocialMediaHandle,
          rcres: rcres
        },
        cache: false,
        success: function() {
          // Success message
          $('#successExecutive').html("<div class='alert alert-success'>");
          $('#successExecutive > .alert-success').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#successExecutive > .alert-success')
            .append("<strong>Thank You for your interest. We will get back to you soon. </strong>");
          $('#successExecutive > .alert-success')
            .append('</div>');
          //clear all fields
          $('#executiveForm').trigger("reset");
          if (rcres.length) {
              grecaptcha.reset();
          }
        },
        error: function(xhr,status,error) {
          if (xhr.responseText != "") {
           var message = JSON.parse(xhr.responseText).Error
            // Fail message
          $('#successExecutive').html("<div class='alert alert-danger'>");
          $('#successExecutive > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#successExecutive > .alert-danger').append($("<strong>").text("Sorry " + message + ", please fill in all required fields"));
          $('#successExecutive > .alert-danger').append('</div>');
          }
          
          //clear all fields
          //$('#applicationForm').trigger("reset");
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
  $('#successExecutive').html('');
});