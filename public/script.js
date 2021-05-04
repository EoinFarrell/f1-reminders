const form = document.querySelector('form');
const thankYouMessage = document.querySelector('#thank-you-message');
form.addEventListener('submit', (e) => {
  e.preventDefault();
  thankYouMessage.classList.add('show');
  setTimeout(() => form.submit(), 2000);
});

// document.addEventListener('DOMContentLoaded', function() { // make sure the DOM is fully loaded before starting anything
//     Notif.setWrapperOptions({ duration: 4000 });
//   }, false);

var successNotif = new Notif("Yay, this message notifies you about the success of whatever!", "success");
var errorNotif = new Notif("Oups, this is a notification you usually hope to not display.", "error");
var confirmedNotif = new Notif("Let's just confirm that whatever happened.", "confirmed");
var defaultNotif = new Notif("This message is just meant to say hi; so \"hi!\"", "default");

var myNotif = new Notif('Hey, what about this very nice notification, with a <a href="#">link</a> and everything?', "default");
document.getElementById('my-button').addEventListener('click', function(e) {
  myNotif.display(3500);
  /* You could also chain everything if it's a one time notification, such as:
   * new Notif('This is a one time notification, so no need to keep it in a JS variable.', "confirmed").display(4000);
   */
});