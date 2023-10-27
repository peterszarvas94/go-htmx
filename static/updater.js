class updaterElement extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.innerHTML = `
      <button id="updater-button" class="cursor-pointer">🔄</button>
    `;

    this.querySelector("#updater-button").addEventListener("click", () => {
      const evt = new Event("ce-update");
      document.dispatchEvent(evt);
    });
  }
}

customElements.define("ce-updater", updaterElement);
