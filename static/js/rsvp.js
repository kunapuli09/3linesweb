$(function() {
  //Check if required fields are filled
    function checkifreqfld() {
        var isFormFilled = true;
        $("#rsvpForm").find(".form-control:visible").each(function() {
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
  $("#rsvpForm input,#rsvpForm textarea").jqBootstrapValidation({
    preventSubmit: true,
    submitError: function($form, event, errors) {
      // additional error messages or events
    },
    submitSuccess: function($form, event) {
      if (checkifreqfld()) {
                event.preventDefault();
      }
      // get values from FORM
      var FullName = $("input#FullName").val();
      var Email = $("input#Email").val();
      var CompanyName = $("input#CompanyName").val();
      var Phone = $("input#Phone").val();
      
      $this = $("#rsvpButton");
      $this.prop("disabled", true); // Disable submit button until AJAX call is complete to prevent duplicate messages
      $.ajax({
        url: "/rsvp",
        type: "POST",
        data: {
          FullName: FullName,
          Email: Email,
          CompanyName: CompanyName,
          Phone: Phone
        },
        cache: false,
        success: function() {
          // Success message
          $('#successRSVP').html("<div class='alert alert-success'>");
          $('#successRSVP > .alert-success').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#successRSVP > .alert-success')
            .append("<strong>Thank You for registering. Zoom Details will be sent in an email soon.</strong>");
          $('#successRSVP > .alert-success')
            .append('</div>');
          //clear all fields
          $('#rsvpForm').trigger("reset");
        },
        error: function(xhr,status,error) {
          if (xhr.responseText != "") {
           var message = JSON.parse(xhr.responseText).Error
            // Fail message
          $('#successRSVP').html("<div class='alert alert-danger'>");
          $('#successRSVP > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#successRSVP > .alert-danger').append($("<strong>").text("Sorry " + message + ", please fill in all required fields"));
          $('#successRSVP > .alert-danger').append('</div>');
          }
          
          //clear all fields
          //$('#rsvpForm').trigger("reset");
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
  $('#successRSVP').html('');
});