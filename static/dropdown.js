class dropdownElement extends HTMLElement {
  constructor() {
    super();
    this.open = false;
  }

  connectedCallback() {
    const items = this.getAttribute("items").split(",");

    const itemsHtml = items.map((item) => `<li>${item}</li>`).join("");

    this.innerHTML = `
      <button id="dropdown-button" class="cursor-pointer">🔽</button>
      <ul id="dropdown-content" class="hidden border absolute bg-white">
        ${itemsHtml}
      </ul>
    `;

    const button = this.querySelector("#dropdown-button");

    button.addEventListener("click", () => {
      this.open = !this.open;
      this.updateOpen();
    });

    document.addEventListener("click", (event) => {
      if (event.target !== button) {
        this.open = false;
        this.updateOpen();
      }
    });
  }

  updateOpen() {
    this.querySelector("#dropdown-content").classList.toggle(
      "hidden",
      !this.open,
    );
    this.querySelector("#dropdown-button").innerText = this.open ? "🔼" : "🔽";
  }
}

customElements.define("ce-dropdown", dropdownElement);
