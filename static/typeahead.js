class typeaheadElement extends HTMLElement {
  constructor() {
    super();
    this.open = false;
  }

  connectedCallback() {
    const items = this.getAttribute("items").split(",");
    const itemsHtml = items
      .map(
        (item) => `
          <li ce-text="${item}" ce-selected="false">
            <button>${item}</button>
          </li>
        `,
      )
      .join("");

    this.innerHTML = `
      <input class="border"/>
      <ul class="hidden">${itemsHtml}</ul>
    `;

    const input = this.querySelector("input");
    const buttons = this.querySelectorAll("button");
    const list = this.querySelectorAll("li");

    // open on focus
    input.addEventListener("focus", () => {
      this.open = true;
      this.updateOpen();
    });

    // change bg on hover
    list.forEach((item) => {
      item.addEventListener("mouseenter", () => {
        this.clearSelection();
        item.setAttribute("ce-selected", true);
        this.updateBg();
      });

      item.addEventListener("mouseleave", () => {
        this.clearSelection();
        this.updateBg();
      });
    });

    // track key presses
    input.addEventListener("keyup", (event) => {
      if (event.key === "Enter") {
        // accept selection
        const selected = this.querySelector(
          "li:not(.hidden)[ce-selected=true]",
        );
        if (selected) {
          input.value = selected.innerText;
        }
        this.updateList();
        return;
      }

      if (event.key === "ArrowDown") {
        // select next
        let index = this.getSelectedIndex();
        const next = this.getNextItem(index);
        this.selectItem(next);
        return;
      }

      if (event.key === "ArrowUp") {
        // select previous
        let index = this.getSelectedIndex();
        const previous = this.getPreviousItem(index);
        this.selectItem(previous);
        return;
      }

      // filter list
      this.updateList();
      const first = this.getItem(0);
      this.selectItem(first);
    });

    // select on click
    buttons.forEach((button) => {
      button.addEventListener("click", () => {
        input.value = button.innerText;
        this.updateList();
        input.focus();
      });
    });

    // close on click outside
    document.addEventListener("click", (event) => {
      const isOutside = !this.contains(event.target);
      if (isOutside) {
        this.open = false;
        this.updateOpen();
      }
    });
  }

  // open/close list
  updateOpen() {
    this.querySelector("ul").classList.toggle("hidden", !this.open);
  }

  // filter list
  updateList() {
    const value = this.querySelector("input").value;
    const list = this.querySelectorAll("li");
    list.forEach((item) => {
      const text = item.getAttribute("ce-text");
      const isMatch = text.toLowerCase().includes(value.toLowerCase());
      item.classList.toggle("hidden", !isMatch);
    });
  }

  // add coloring on arrow navigation
  updateBg() {
    const list = this.querySelectorAll("li");
    list.forEach((item) => {
      item.classList.toggle(
        "bg-gray-200",
        item.getAttribute("ce-selected") === "true",
      );
    });
  }

  // clear selection
  clearSelection() {
    const list = this.querySelectorAll("li");
    list.forEach((item) => {
      item.setAttribute("ce-selected", false);
    });
  }

  // get selected item
  getSelectedIndex() {
    const list = this.querySelectorAll("li:not(.hidden)");
    const length = list.length;
    for (let i = 0; i < length; i++) {
      const item = list[i];
      const isSelected = item.getAttribute("ce-selected") === "true";
      if (isSelected) {
        return i;
      }
    }
    return -1;
  }

  // get item
  getItem(index) {
    const list = this.querySelectorAll("li:not(.hidden)");
    return list[index];
  }

  // get next item
  getNextItem(index) {
    const list = this.querySelectorAll("li:not(.hidden)");
    const length = list.length;
    const nextIndex = (index + 1) % length;
    return list[nextIndex];
  }

  // get previous item
  getPreviousItem(index) {
    const list = this.querySelectorAll("li:not(.hidden)");
    const length = list.length;
    const previousIndex = (index - 1 + length) % length;
    return list[previousIndex];
  }

  // select an item
  selectItem(item) {
    this.clearSelection();
    item.setAttribute("ce-selected", true);
    this.updateBg();
  }
}

customElements.define("ce-typeahead", typeaheadElement);
