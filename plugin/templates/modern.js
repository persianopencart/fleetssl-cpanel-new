/*
 * FleetSSL cPanel - lightweight UI helpers.
 * Vanilla JavaScript, no external dependencies (no jQuery, no CDN).
 */

/* Disables submit buttons once a form is submitted to prevent double posts. */
function LetsEncrypt_DisableButtons() {
  var buttons = document.querySelectorAll(
    ".fleetssl-ui input[type='submit'], .fleetssl-ui button[type='submit']");
  for (var i = 0; i < buttons.length; i++) {
    buttons[i].disabled = true;
    if (buttons[i].tagName === "INPUT") {
      buttons[i].value = "Please wait...";
    }
  }
  document.body.style.cursor = "wait";
  return true;
}

/*
 * On the issuance page, wildcard names can only be validated with dns-01, so
 * the wildcard checkboxes are enabled only while the dns-01 method is chosen.
 */
document.addEventListener("DOMContentLoaded", function () {
  var radios = document.querySelectorAll("input[name='challenge_method']");
  if (!radios.length) {
    return;
  }

  function applyMethod() {
    var selected = document.querySelector("input[name='challenge_method']:checked");
    if (!selected) {
      return;
    }
    var wildcards = document.querySelectorAll(".wildcard-checkbox");
    for (var i = 0; i < wildcards.length; i++) {
      if (selected.value === "dns-01") {
        wildcards[i].disabled = false;
      } else {
        wildcards[i].checked = false;
        wildcards[i].disabled = true;
      }
    }
  }

  for (var i = 0; i < radios.length; i++) {
    radios[i].addEventListener("change", applyMethod);
  }
  applyMethod();
});
