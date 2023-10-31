// signal implementation by Awesome
// https://www.youtube.com/watch?v=bUy4xiJ05KY

let RUNNING = null;

function runAndExtractDependencies(task) {
  RUNNING = task;
  task.execute();
  RUNNING = null;
}

class Signal {
  constructor(initValue) {
    this._value = initValue;
    this.dependencies = [];
  }

  get value() {
    if (RUNNING) this.dependencies.push(RUNNING);
    return this._value;
  }

  set value(newValue) {
    if (this._value === newValue) return;
    this._value = newValue;
    this.notify();
  }

  notify() {
    for (const dep of this.dependencies) {
      dep.update();
    }
  }
}

class Computed {
  constructor(compute) {
    this.compute = compute;
    this.isStale = true;
    runAndExtractDependencies(this);
  }

  get value() {
    if (this.isStale) {
      this._value = this.compute();
      this.isStale = false;
    }
    return this._value;
  }

  execute() {
    this.compute();
  }

  update() {
    this.isStale = true;
  }
}

class Effect {
  constructor(effect) {
    this.effect = effect;
    runAndExtractDependencies(this);
  }

  execute() {
    this.effect();
  }

  update() {
    this.execute();
  }
}

export function signal(s) {
  return new Signal(s);
}

export function computed(c) {
  return new Computed(c);
}

export function effect(e) {
  return new Effect(e);
}
