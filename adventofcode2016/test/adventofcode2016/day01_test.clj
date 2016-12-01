(ns adventofcode2016.day01-test
  (:require [clojure.test :refer :all]
            [adventofcode2016.day01 :refer :all]))

(deftest parse-test
  (testing "Parse two commands"
    (is (= 2 (count (parse "R2, L3")))))
  (testing "Verify direction, first item"
    (is (= 1 (:direction (first (parse "R2, L3"))))))
  (testing "Verify blocks, first item"
    (is (= 2 (:blocks (first (parse "R2, L3"))))))
  (testing "Verify direction, second item"
    (is (= -1 (:direction (last (parse "R2, L3"))))))
  (testing "Verify blocks, second item"
    (is (= 3 (:blocks (last (parse "R2, L3")))))))

(deftest get-distance-test
  (testing "Positive numbers"
    (is (= 4 (get-distance {:row 2 :col 2}))))
  (testing "Negative numbers"
    (is (= 10 (get-distance {:row -2 :col -8}))))
  (testing "Mixed numbers"
    (is (= 4 (get-distance {:row -2 :col 2})))))

(deftest walk-test
  (testing "Basic move"
    (is (= { :row 5 :col 2 } (walk (parse "R2, L5")))))
  (testing "Almost a circle"
    (is (= { :row -2 :col 0 } (walk (parse "R2, R2, R2")))))
  (testing "Circle twice"
    (is (= { :row 0 :col 0 } (walk (parse "R2, R2, R2, R2, L2, L2, L2, L2")))))
  (testing "Another example"
    (is (= { :row 2 :col 10 } (walk (parse "R5, L5, R5, R3"))))))

(deftest solve-test
  (testing "Basic move"
    (is (= 7 (solve "R2, L5"))))
  (testing "Almost a circle"
    (is (= 2 (solve "R2, R2, R2"))))
  (testing "Circle twice"
    (is (= 0 (solve "R2, R2, R2, R2, L2, L2, L2, L2"))))
  (testing "Another example"
    (is (= 12 (solve "R5, L5, R5, R3")))))
