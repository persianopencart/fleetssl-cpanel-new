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

/*
 * Tab switching for the WHM page. Clicking a .nav-tabs link reveals its
 * matching .tab-pane. Replaces the old Bootstrap/jQuery tab plugin.
 */
document.addEventListener("DOMContentLoaded", function () {
  var links = document.querySelectorAll(".nav-tabs a");
  for (var i = 0; i < links.length; i++) {
    links[i].addEventListener("click", function (e) {
      e.preventDefault();
      var target = document.querySelector(this.getAttribute("href"));
      if (!target) {
        return;
      }
      var tabs = document.querySelectorAll(".nav-tabs li");
      for (var j = 0; j < tabs.length; j++) {
        tabs[j].classList.remove("active");
      }
      var panes = document.querySelectorAll(".tab-content > .tab-pane");
      for (var j = 0; j < panes.length; j++) {
        panes[j].classList.remove("active");
      }
      this.parentNode.classList.add("active");
      target.classList.add("active");
    });
  }
});
