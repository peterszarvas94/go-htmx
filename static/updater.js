class updaterElement extends HTMLElement {
  constructor() {
    super();
    this.eventName = "";
  }

  connectedCallback() {
    this.innerHTML = `
      <button id="updater-button" class="cursor-pointer">🔄</button>
    `;

    this.eventName = this.getAttribute("ce-event");

    this.querySelector("#updater-button").addEventListener("click", () => {
      if (!this.eventName) {
        console.error("No event name specified");
        return;
      }
      const evt = new Event(this.eventName);
      document.dispatchEvent(evt);
    });

  }
}

customElements.define("ce-updater", updaterElement);
