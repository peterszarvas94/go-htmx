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
          <li class="w-full" ce-text="${item}" ce-selected="false">
            <button class="w-full text-left">${item}</button>
          </li>
        `,
      )
      .join("");

    this.innerHTML = `
      <input class="border border-black rounded-sm w-40"/>
      <ul class="hidden w-40 border absolute bg-white">${itemsHtml}</ul>
    `;

    const input = this.querySelector("input");
    const buttons = this.querySelectorAll("button");
    const list = this.querySelectorAll("li");

    // change bg on hover
    list.forEach((item) => {
      item.addEventListener("mouseenter", () => {
        this.selectItem(item);
        this.updateStyle();
      });

      item.addEventListener("mouseleave", () => {
        this.clearSelection();
        this.updateStyle();
      });
    });

    input.addEventListener("focus", () => {
      this.updateOpen(true);
    });

    input.addEventListener("click", () => {
      this.updateOpen(true);
    });

    // track key presses
    input.addEventListener("keyup", (event) => {
      if (event.key === "Escape") {
        // close list
        this.updateOpen(false);
        return;
      }

      if (event.key === "Enter") {
        if (!this.open) {
          // open list
          this.updateOpen(true);
          return;
        }

        // accept selection
        const selected = this.querySelector(
          'li:not(.hidden)[ce-selected="true"]',
        );
        if (!selected) {
          return;
        }
        input.value = selected.innerText;
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
      if (!first) {
        return;
      }
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
        this.updateOpen(false);
      }
    });

    // close on tab
    document.addEventListener("keydown", (event) => {
      if (event.key === "Tab") {
        this.updateOpen(false);
      }
    });
  }

  // open/close list
  updateOpen(bool) {
    this.open = bool;
    this.querySelector("ul").classList.toggle("hidden", !bool);
    const first = this.getItem(0);
    if (!first) {
      return;
    }
    this.selectItem(first);
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
  updateStyle() {
    const list = this.querySelectorAll("li");
    const classes = ["bg-blue-600", "text-white"];
    list.forEach((item) => {
      classes.forEach((className) => {
        item.classList.toggle(
          className,
          item.getAttribute("ce-selected") === "true",
        );
      });
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
    this.updateStyle();
  }
}

customElements.define("ce-typeahead", typeaheadElement);
