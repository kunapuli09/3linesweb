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
      var Title = $("input#Title").val();
      /* declare an checkbox array */
      /* we join the array separated by the comma */
      var industriesArray = [];
      var industries;
      /* look for all checkboes that have a class 'Industries' attached to it and check if it was checked */
      $(".form-check-input:checked").each(function() {
        industriesArray.push($(this).val());
      });
      industries = industriesArray.join(',') ;
      // if(industries.length > 0){
      //   alert("You have selected " + industries); 
      // }else{
      //   alert("Please at least check one of the checkbox"); 
      // }
      /* declare an checkbox array */
      /* we join the array separated by the comma */
      var locationsArray = [];
      var locations;
      
      /* look for all checkboes that have a class 'chk' attached to it and check if it was checked */
      $(".form-check-input1:checked").each(function() {
        locationsArray.push($(this).val());
      });
      locations = locationsArray.join(',') ;
      // if(locations.length > 0){
      //   alert("You have selected " + locations); 
      // }else{
      //   alert("Please at least check one of the checkbox"); 
      // }
      var Industries = industries;
      var Locations = locations;
      var CapitalRaised = $("input#CapitalRaised").val();
      var Comments = $("textarea#Comments").val();
      
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
          Title: Title,
          Industries: Industries,
          Locations: Locations,
          CapitalRaised: CapitalRaised,
          Comments: Comments
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
        error: function() {
          // Fail message
          $('#successApplication').html("<div class='alert alert-danger'>");
          $('#successApplication > .alert-danger').html("<button type='button' class='close' data-dismiss='alert' aria-hidden='true'>&times;")
            .append("</button>");
          $('#successApplication > .alert-danger').append($("<strong>").text("Sorry " + FirstName + ", please fill in all required fields"));
          $('#successApplication > .alert-danger').append('</div>');
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