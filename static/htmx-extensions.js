// "auth" extension to render error messages on forms
htmx.defineExtension("auth", {
  onEvent: function (name, evt) {
    if (name == "htmx:beforeOnLoad") {
      const status = evt.detail.xhr.status;
      const statusUnauthorized = status === 401;
      const statusForbidden = status === 403;
      const statusConflict = status === 409;
      const displayErrorMessage =
        statusUnauthorized || statusForbidden || statusConflict;
      if (displayErrorMessage) {
        evt.detail.shouldSwap = true;
      }
    }
  },
});

// "delete" extension to delete elements
htmx.defineExtension("delete", {
  onEvent: function (_name, evt) {
    const status = evt?.detail?.xhr?.status;
    if (status === 204) {
      evt.detail.shouldSwap = true;
    }
  },
});

// "error" extension to prevent swap on 500 error
htmx.defineExtension("error", {
  onEvent: function (_name, evt) {
    const erroStatus = evt?.detail?.xhr?.status;
    if (erroStatus === 500) {
      evt.detail.shouldSwap = false;
    }
  },
});

let accessToken = undefined;

// "protected" extension to include authorization header
htmx.defineExtension("protected", {
  onEvent: function (name, evt) {
    if (name === "htmx:beforeRequest") {
      if (accessToken) {
        evt.detail.xhr.setRequestHeader(
          "Authorization",
          "Bearer " + accessToken,
        );
      }
    }
  },
});

// "signout" extension to clear access token
htmx.defineExtension("signout", {
  onEvent: function (name, evt) {
    if (name === "htmx:afterRequest") {
      if (evt.detail.xhr.status == 200) {
        accessToken = undefined;
      }
    }
  },
});

// "signin" extension to set access token
htmx.defineExtension("signin", {
  onEvent: function (name, evt) {
    if (name === "htmx:afterRequest" && evt.detail.xhr.status == 200) {
      auth = evt.detail.xhr.getResponseHeader("Authorization");
      accessToken = auth.split(" ")[1];
    }
  },
});

// "token" extension to get access token from server
htmx.defineExtension("token", {
  onEvent: async function (name, evt) {
    if (name === "htmx:load" && evt.target === document.body) {
      if (accessToken) {
        return;
      }

      try {
        const res = await fetch("/refresh", { method: "POST" });
        auth = res.headers.get("Authorization");
        if (auth) {
          accessToken = auth.split(" ")[1];
          await htmx.ajax("GET", window.location.pathname, {
            headers: {
              Authorization: "Bearer " + accessToken,
            },
            target: "body",
          });
        }
      } catch {}
    }
  },
});

/*
Template extension
Example:
hx-get="/some/fake/url"
hx-ext="template"
hx-template="#template-id"
hx-template-text="variable-name" (optional, variable come from url params)
hx-template-run="js-code" (optional, js code to run before template is rendered)
*/

/*
htmx.defineExtension("template", {
  onEvent: function (name, evt) {
    if (name === "htmx:beforeSend") {
      document.dispatchEvent(new Event("htmx:xhr:abort"));
    }
  },

  transformResponse: function (text, xhr, elt) {
    // get template
    const query = elt.getAttribute("hx-template");
    const template = document.querySelector(query);
    const content = template.content;

    // set values for window
    const url = new URL(xhr.responseURL);
    const params = url.searchParams;
    for (const [key, value] of params) {
      // try to parse to number
      const num = Number(value);
      if (!isNaN(num)) {
        window[key] = num;
        continue;
      }
      window[key] = value;
    }

    // set text for elements inside template
    const textElts = content.querySelectorAll("[hx-template-text]");
    textElts.forEach((textElt) => {
      const name = textElt.getAttribute("hx-template-text");
      const value = eval(name);
      textElt.innerText = value;
    });

    // check if run attribute exists
    const run = elt.getAttribute("hx-template-run");
    if (run) {
      eval(run);
    }

    // set values for imputs inside template
    const inputs = content.querySelectorAll("input");
    inputs.forEach((input) => {
      const name = input.getAttribute("name");
      const value = window[name];
      input.setAttribute("value", value);
    });

    const serializer = new XMLSerializer();
    const str = serializer.serializeToString(content);
    return str;
  },
});
*/
