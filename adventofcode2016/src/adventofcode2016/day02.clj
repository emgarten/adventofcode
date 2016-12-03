(ns adventofcode2016.day02
  (:use [clojure.pprint])
  (:require [clojure.string :as str]))

(defn to-key-id
  "x -> Base3 - 1"
  [x]
  (Integer. (Integer/toString (- x 1) 3)))

(defn from-key-id
  "x -> Base10 + 1"
  [x]
  (+ 1 (+ (mod x 10) (* 3 (int (Math/floor (/ x 10)))))))

(defn letter-fun
  "Get letter function"
  [letter]
  (case (str letter)
    "L" #(- % 1)
    "R" #(+ % 1)
    "U" #(- % 10)
    "D" #(+ % 10)
    #(identity %)))

(defn apply-letter-or-skip
  [key-id letter]
  (let [next-id ((letter-fun letter) key-id)]
    (if
      (and
       (> next-id -1)
       (< next-id 23)
       (> 3 (mod next-id 10)))
      ; valid
      next-id
      ; noop
      key-id)))

(defn process-single-key
  [[letter & other-letters] key-id]
  (cond
    (nil? letter) key-id
    :else (recur
            other-letters
            (apply-letter-or-skip key-id letter))))

(defn process-keys
  [[line & other-lines] key-ids]
  (cond
    (nil? line) key-ids
    :else (recur other-lines (conj key-ids (process-single-key line (last key-ids))))))

(defn get-lines
  [input-data]
  (filter #(> (count %) 0)
    (map #(str/trim %)
      (str/split-lines input-data))))

(defn solve-part-1
  "Read and process day2 input to find the keypad combination."
  [input-data]
  (map #(from-key-id %)
    (rest
     (process-keys
      (get-lines input-data)
      [11]))))
