function formSubmit(event) {
  var url = "https://phrhyp7dx2.execute-api.eu-west-1.amazonaws.com/Production";
  var request = new XMLHttpRequest();
  request.open('POST', url, true);
  request.onload = function() { // request successful
  // we can use server response to our request now
    console.log(request.responseText);
    new Notif("Thank you for signing up and enjoy the next race.", "success").display(3500);
  };

  request.onerror = function() {
    // request failed
    new Notif("Error submitting form, sorry, maybe try again.", "error").display(3500);
  };

  request.send(new FormData(event.target)); // create FormData from form that triggered event
  event.preventDefault();
}


document.querySelector('form').addEventListener("submit", formSubmit);