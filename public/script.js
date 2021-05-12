async function formSubmit(event) {
  event.preventDefault();
  
  let user = {
    phone: document.querySelector('#phone').value,
    email: document.querySelector('#email').value,
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
  };
  
  let response = await fetch("https://phrhyp7dx2.execute-api.eu-west-1.amazonaws.com/Production", {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(user)
  });

  if (response.ok) { // if HTTP-status is 200-299
    new Notif("Thank you for signing up and enjoy the next race.", "success").display(3500);
  } else {
    new Notif("Error submitting form, sorry, maybe try again. Debug info: " + response.status, "error").display(3500);
  }
}

document.querySelector('form').addEventListener("submit", formSubmit);