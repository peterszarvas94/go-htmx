import { signal, effect, computed } from "/static/signal.js";

class Signaler extends HTMLElement {
  constructor() {
    super();
    this.count = signal(0);
    this.doubleCount = computed(() => this.count.value * 2);
  }

  connectedCallback() {
    this.innerHTML = `
      <button id="increment">➕</button>
      <div id="container">${this.count.value} x 2 = ${this.doubleCount.value}</div>
    `;

    const container = this.querySelector("#container");
    const button = this.querySelector("#increment");

    effect(() => {
      container.innerText = `${this.count.value} x 2 = ${this.doubleCount.value}`;
    });

    button.addEventListener("click", () => {
      this.count.value++;
    });
  }
}

customElements.define("ce-signaler", Signaler);
