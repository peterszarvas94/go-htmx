class counterElement extends HTMLElement {
  constructor() {
    super();
    this.count = 0;
  }

  connectedCallback() {
    const start = parseInt(this.getAttribute("start"));
    if (!isNaN(start)) {
      this.count = start;
    }

    this.innerHTML = `
      <div class="flex gap-2">
        <button id="decrement">➖</button>
        <div id="counter">${this.count}</div>
        <button id="increment">➕</button>
      </div>
    `;

    this.querySelector("#increment").addEventListener("click", () => {
      this.count++;
      this.updateCount();
    });

    this.querySelector("#decrement").addEventListener("click", () => {
      this.count--;
      this.updateCount();
    });
  }

  updateCount() {
    this.querySelector("#counter").innerHTML = this.count;
  }
}

customElements.define("ce-counter", counterElement);
