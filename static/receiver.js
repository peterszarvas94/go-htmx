class receiverElement extends HTMLElement {
  constructor() {
    super();
    this.count = 0;
  }

  connectedCallback() {
    this.innerHTML = `
      <div id="receiver">Event received: ${this.count}</div>
    `;

    document.addEventListener("ce-update", () => {
      this.count++;
      this.querySelector("#receiver").innerText = `Event received: ${this.count}`;
    });
  }
}

customElements.define("ce-receiver", receiverElement);
