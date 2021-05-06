function formSubmit(event) {
  var url = "https://phrhyp7dx2.execute-api.eu-west-1.amazonaws.com/Production";
  
  var request = new XMLHttpRequest();
  request.open('POST', url, true);
  request.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
  
  request.onload = function() {
    console.log(request.responseText);
    new Notif("Thank you for signing up and enjoy the next race.", "success").display(3500);
  };

  request.onerror = function() {
    new Notif("Error submitting form, sorry, maybe try again.", "error").display(3500);
  };

  request.send(new FormData(event.target));
  event.preventDefault();
}

document.querySelector('form').addEventListener("submit", formSubmit);