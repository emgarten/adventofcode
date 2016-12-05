(ns adventofcode2016.day04
  (:use [adventofcode2016.common])
  (:use [clojure.pprint])
  (:require [clojure.string :as str]))

(defrecord Room [data sector letters checksum])
(defrecord LetterGroup [letter count])

(defn parse-line
  [line]
  (let [fb (str/index-of line "[")
        sb (str/index-of line "]")
        data (subs line 0 fb)
        checksum (subs line (inc fb) sb)
        letters (filter #(Character/isLowerCase %) data)
        sector (Integer. (str/join (filter #(Character/isDigit %) data)))]
    ;(pprint (str data " " sector " " letters " " checksum))
    (->Room data sector letters checksum)))

(defn get-rooms
  [lines]
  (map #(parse-line %) lines))

(defn compare-groups
  [x y]
  ; Compare on count
  (let [c (compare (:count y) (:count x))]
    (if (not= c 0)
      ; count is not equal
      c
      ; counts are equal, compare on letter
      (compare (:letter x) (:letter y)))))

(defn create-checksum
  [room]
  (let [groups (group-by identity (:letters room))
        counts (map #(->LetterGroup (key %) (-> % val count)) groups)]
    (str/join
     (map :letter
       (take 5
             (sort compare-groups counts))))))

(defn valid-checksum?
  [room]
  (let [checksum (create-checksum room)]
    ;(pprint (str "checksum: " checksum " match: " (:checksum room)))
    (= checksum (:checksum room))))

(defn solve-part-1
  [input-data]
  (apply +
   (map :sector
      (filter #(valid-checksum? %) (get-rooms (get-lines input-data))))))
