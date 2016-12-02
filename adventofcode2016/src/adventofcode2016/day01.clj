(ns adventofcode2016.day01
  (:use [clojure.pprint])
  (:require [clojure.string :as str]))

(defn add-directions
  "Find the resulting position after adding two positions together."
  [a b]
  {:col (+ (:col a) (:col b))
   :row (+ (:row a) (:row b))})

(defn parse-direction-value
  "Convert R -> 1, L -> -1"
  [entry]
  (cond
    (= "R" (subs entry 0 1)) 1
    :else -1))

(defn parse-entry
  "Parse a single direction, example: R2"
  [entry]
  {:direction (parse-direction-value entry)
   :blocks (Integer. (subs entry 1))})

(defn parse
  "Parse a comma delimited list of directions, example: R2, L3"
  [input-data]
  (map parse-entry
       (map #(str/trim %)
            (str/split input-data #","))))

(defn get-distance
  "Distince from 0,0 on a grid."
  [position]
  (+ (Math/abs (:row position)) (Math/abs (:col position))))

(defn walk-translate-entry
  "Translate a relative command to an absolute command"
  [facing blocks]
  (cond
    ; Up, North
    (= 0 facing) {:col 0 :row blocks}
    ; Right, East
    (= 1 facing) {:col blocks :row 0}
    ; Down, South
    (= 2 facing) {:col 0 :row (* -1 blocks)}
    ; Left, West
    (= 3 facing) {:col (* -1 blocks) :row 0}
    ; Invalid data
    :else (throw (Exception. "Invalid facing"))))

(defn walk-translate
  "Translate relative commands to absolute commands"
  ([commands]
   (walk-translate 0 commands []))
  ([facing [command & other-commands] translated]
   (cond
     ; Done
     (nil? command) translated
     :else
       ; Find the new direction and move towards it
       ; There are only four possible directions, so here we will
       ; turn in either a positive or negative way and take mod 4
       ; to get the direction of travel.
     (let [{direction :direction
            blocks :blocks} command
           new-facing (mod (+ facing direction) 4)
           translated-entry (walk-translate-entry new-facing blocks)]
       (recur
        new-facing
        other-commands
        (conj translated translated-entry))))))

(def start-position
  {:row 0 :col 0})

(defn inc-row
  [{row :row
    col :col}]
  {:row (inc row) :col col})

(defn dec-row
  [{row :row
    col :col}]
  {:row (dec row) :col col})

(defn inc-col
  [{row :row
    col :col}]
  {:row row :col (inc col)})

(defn dec-col
  [{row :row
    col :col}]
  {:row row :col (dec col)})

(defn walk-translate-to-final
  "Translate all positions to their final grid position."
  ([commands]
   (walk-translate-to-final start-position commands [start-position]))
  ([position [command & other-commands] translated-to-final]
   (cond
     (nil? command) translated-to-final
     :else
     (let [new-position (add-directions position command)]
       (recur
        new-position
        other-commands
        (conj translated-to-final new-position))))))

(defn get-incremental-positions
  "Get all points between positions."
  ([current destination positions]
   (cond
     (nil? destination) [current]
     (= current destination) positions
     :else (let [{current-row :row
                  current-col :col} current
                 {dest-row :row
                  dest-col :col} destination]
             (cond
               (< current-row dest-row) (recur (inc-row current) destination (conj positions current))
               (> current-row dest-row) (recur (dec-row current) destination (conj positions current))
               (< current-col dest-col) (recur (inc-col current) destination (conj positions current))
               (> current-col dest-col) (recur (dec-col current) destination (conj positions current)))))))

(defn walk-incremental-positions
  "Get all points hit on the walk."
  ([positions]
   (walk-incremental-positions positions '()))
  ([[position & other-positions] all-positions]
   (cond
    (nil? position) all-positions
    :else
     (recur other-positions
      (flatten
       (concat
        all-positions
        (get-incremental-positions position (first other-positions) [])))))))

(defn walk
  "Process commands and return the end position."
  [commands]
  (-> commands
      walk-translate
      walk-translate-to-final
      last))

(defn solve-part-1
  "Read and process day1 input to find the end distance."
  [input-data]
  (-> input-data
      parse
      walk
      get-distance))

(defn find-dupe
  [positions]
  (loop [seen '()
         search positions]
    (let [current (first search)]
      (cond
        (nil? current) nil
        (some #(= current %) seen) current
        :else (recur (conj seen current) (rest search))))))

(defn solve-part-2
  "Read and process day1 input to find the first duplicate entry."
  [input-data]
  (-> input-data
      parse
      walk-translate
      walk-translate-to-final
      walk-incremental-positions
      find-dupe
      get-distance))
