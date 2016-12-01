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

(defn walk
  "Process commands and return the end position."
  [commands]
  ; Translate commands from relative to absolute
  ; Then add the absolute movements together to
  ; get the final position.
  (reduce add-directions
    (walk-translate commands)))

(defn solve
  "Read and process day1 input."
  [input-data]
  (-> input-data
      parse
      walk
      get-distance))
