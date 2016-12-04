(ns adventofcode2016.day02-test
  (:require [clojure.test :refer :all]
            [adventofcode2016.day02 :refer :all]))

(deftest solve-part-1-test
  (testing "Single line"
    (is (= '("1") (solve-part-1 "ULL"))))
  (testing "Four lines"
    (is (= '("1" "9" "8" "5") (solve-part-1 "ULL
                               RRDDD
                               LURDL
                               UUUUD \n "))))
  (testing "Complex"
    (is (= '("9" "5") (solve-part-1 "DRRRRRRRRRRRRRUUUUUUUDUDUDUDUDULLDDRR\nULLR")))))

(deftest solve-part-2-test
  (testing "Single line"
    (is (= '("5") (solve-part-2 "ULL"))))
  (testing "Four lines"
    (is (= '("5", "D", "B", "3") (solve-part-2 "ULL
                                RRDDD
                               LURDL
                               UUUUD \n ")))))

(deftest get-key-from-keypad-by-id-test
  (testing "1"
    (is (= "1" (:id (get-key-from-keypad keypad-1 "1"))))))

(deftest get-key-from-keypad-by-row-col-test
  (testing "1"
     (is (= "1" (:id (get-key-from-keypad keypad-1 0 0)))))
  (testing "2"
    (is (= "2" (:id (get-key-from-keypad keypad-1 0 1)))))
  (testing "9"
    (is (= "9" (:id (get-key-from-keypad keypad-1 2 2))))))
