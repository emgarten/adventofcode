(ns adventofcode2016.day02
  (:use [clojure.pprint])
  (:require [clojure.string :as str]))

; Key on a keypad
(defrecord Key [id row col])

; Part 1
(def keypad-1
  [
    (->Key "1" 0 0)
    (->Key "2" 0 1)
    (->Key "3" 0 2)
    (->Key "4" 1 0)
    (->Key "5" 1 1)
    (->Key "6" 1 2)
    (->Key "7" 2 0)
    (->Key "8" 2 1)
    (->Key "9" 2 2)])

(defn match-key-id
  [key key-id]
  (= key-id (:id key)))

(defn get-key-from-keypad
  ([keypad key-id]
   (let [result (first (filter #(match-key-id % key-id) keypad))]
     result))
  ([keypad row col]
   (let [result
         (first
          (filter
             #(and (= row (:row %)) (= col (:col %)))
           keypad))]
     result)))

(defn get-next-key-id
  [keypad key-id letter]
  (let [current-key (get-key-from-keypad keypad key-id)
        next-key
          (case (str letter)
            "L" (get-key-from-keypad keypad (:row current-key) (dec (:col current-key)))
            "R" (get-key-from-keypad keypad (:row current-key) (inc (:col current-key)))
            "U" (get-key-from-keypad keypad (dec (:row current-key)) (:col current-key))
            "D" (get-key-from-keypad keypad (inc (:row current-key)) (:col current-key))
            (get-key-from-keypad keypad key-id))]
      (cond
        (nil? next-key) key-id
        :else (:id next-key))))

(defn process-single-key
  [keypad [letter & other-letters] key-id]
  (cond
    (nil? letter) key-id
    :else (recur
            keypad
            other-letters
            (get-next-key-id keypad key-id letter))))

(defn process-keys
  [keypad [line & other-lines] key-ids]
  (cond
    (nil? line) key-ids
    :else (recur keypad other-lines (conj key-ids (process-single-key keypad line (last key-ids))))))

(defn get-lines
  [input-data]
  (filter #(> (count %) 0)
    (map #(str/trim %)
      (str/split-lines input-data))))

(defn solve-part-1
  "Read and process day2 input to find the keypad combination."
  [input-data]
  (rest
   (process-keys
    keypad-1
    (get-lines input-data)
    ["5"])))
