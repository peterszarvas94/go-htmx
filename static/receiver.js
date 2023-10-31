class receiverElement extends HTMLElement {
  constructor() {
    super();
    this.count = 0;
    this.eventName = "";
  }

  connectedCallback() {
    this.innerHTML = `
      <div id="receiver">Event received: ${this.count}</div>
    `;

    this.eventName = this.getAttribute("ce-event");

    document.addEventListener(this.eventName, () => {
      this.count++;
      this.querySelector("#receiver").innerText = `Event received: ${this.count}`;
    });
  }
}

customElements.define("ce-receiver", receiverElement);
