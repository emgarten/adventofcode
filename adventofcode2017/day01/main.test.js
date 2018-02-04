const sumMatch = require('./main');

test('1122', () => {
  expect(sumMatch("1122")).toBe(3);
});

test('1111', () => {
  expect(sumMatch("1111")).toBe(4);
});

test('1234', () => {
  expect(sumMatch("1234")).toBe(0);
});

test('91212129', () => {
  expect(sumMatch("91212129")).toBe(9);
});