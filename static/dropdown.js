class dropdownElement extends HTMLElement {
  constructor() {
    super();
    this.open = false;
  }

  connectedCallback() {
    const attrOpen = this.getAttribute("open");
    if (attrOpen === "true") {
      this.open = true;
    }

    this.innerHTML = `
      <button id="dropdown-button" class="cursor-pointer">🔽</button>
      <ul id="dropdown-content" class="hidden">
        <li>Item 1</li>
        <li>Item 2</li>
        <li>Item 3</li>
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
    this.querySelector("#dropdown-content").classList.toggle("hidden", !this.open);
    this.querySelector("#dropdown-button").innerText = this.open ? "🔼" : "🔽";
  }
}

customElements.define("ce-dropdown", dropdownElement);
