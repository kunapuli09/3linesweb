$(function() {

  $("#applicationForm input,#applicationForm textarea").jqBootstrapValidation({
    preventSubmit: true,
    submitError: function($form, event, errors) {
      // additional error messages or events
    },
    submitSuccess: function($form, event) {
      event.preventDefault(); // prevent default submit behaviour
      // get values from FORM
      var FirstName = $("input#FirstName").val();
      var LastName = $("input#LastName").val();
      var Email = $("input#Email").val();
      var CompanyName = $("input#CompanyName").val();
      var Website = $("input#Website").val();
      var Phone = $("input#Phone").val();
      var Referrer = $("input#Referrer").val();
      /* declare an checkbox array */
      /* we join the array separated by the comma */
      var industriesArray = [];
      var industries;
      /* look for all checkboes that have a class 'Industries' attached to it and check if it was checked */
      $(".form-check-input:checked").each(function() {
        industriesArray.push($(this).val());
      });
      industries = industriesArray.join(',') ;
      /* declare an checkbox array */
      /* we join the array separated by the comma */
      var locationsArray = [];
      var locations;
      locationsArray.push($( "#Program option:selected" ).val());
      locations = locationsArray.join(',') ;

       
      var revenueArray = [];
      var revenue;
      
      /* look for all checkboes that have a class 'chk' attached to it and check if it was checked */
      $(".form-check-input2:checked").each(function() {
        revenueArray.push($(this).val());
      });
      revenue = revenueArray.join(',') ;
      var Industries = industries;
      var Locations = locations;
      var Revenue = revenue;
      var CapitalRaised = $("input#CapitalRaised").val();
      var Comments = $("textarea#Comments").val();
      var ElevatorPitch = $("textarea#ElevatorPitch").val();
      var firstName = FirstName; // For Success/Failure Message
      // Check for white space in name for Success/Fail message
      if (firstName.indexOf(' ') >= 0) {
        firstName = name.split(' ').slice(0, -1).join(' ');
      }
      $this = $("#sendApplicationButton");
      $this.prop("disabled", true); // Disable submit button until AJAX call is complete to prevent duplicate messages
      $.ajax({
        url: "/application",
        type: "POST",
        data: {
          FirstName: FirstName,
          LastName: LastName,
          Email: Email,
          CompanyName: CompanyName,
          Website: Website,
          Phone: Phone,
          Referrer: Referrer,
          Industries: Industries,
          Locations: Locations,
          Revenue: Revenue,
          CapitalRaised: CapitalRaised,
          Comments: Comments,
          ElevatorPitch: ElevatorPitch
        },
        cache: false,
        success: function() {
          // Success message
          $('#successApplication').html("<div class='alert alert-success'>");
          $('#successApplication > .alert-success').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#successApplication > .alert-success')
            .append("<strong>Thank You for applying. We will get back to you soon. </strong>");
          $('#successApplication > .alert-success')
            .append('</div>');
          //clear all fields
          $('#applicationForm').trigger("reset");
        },
        error: function(xhr,status,error) {
          if (xhr.responseText != "") {
           var message = JSON.parse(xhr.responseText).Error
            // Fail message
          $('#successApplication').html("<div class='alert alert-danger'>");
          $('#successApplication > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#successApplication > .alert-danger').append($("<strong>").text("Sorry " + message + ", please fill in all required fields"));
          $('#successApplication > .alert-danger').append('</div>');
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
  $('#successApplication').html('');
});