//const sumMatch = require('./main');
const main = require('./main');
const sumMatch = main.sumMatch;
const sumMatchHalf = main.sumMatchHalf;

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

test('1212', () => {
  expect(sumMatchHalf("1212")).toBe(6);
});

test('1221', () => {
  expect(sumMatchHalf("1221")).toBe(0);
});

test('123425', () => {
  expect(sumMatchHalf("123425")).toBe(4);
});

test('123123', () => {
  expect(sumMatchHalf("123123")).toBe(12);
});

test('12131415', () => {
  expect(sumMatchHalf("12131415")).toBe(4);
});