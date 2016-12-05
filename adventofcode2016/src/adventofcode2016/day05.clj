(ns adventofcode2016.day05
  (:use [adventofcode2016.common])
  (:use [clojure.pprint])
  (:require [clojure.string :as str]))

; clojure md5 https://gist.github.com/jizhang/4325757,
(import 'java.security.MessageDigest
        'java.math.BigInteger)

(def md5Instance (MessageDigest/getInstance "MD5"))
(def md5Size (* 2 (.getDigestLength md5Instance)))

(defn md5 [s]
  (let [algorithm md5Instance
        size md5Size
        raw (.digest algorithm (.getBytes s))
        sig (.toString (BigInteger. 1 raw) 16)
        padding (apply str (repeat (- size (count sig)) "0"))]
    (str padding sig)))

(def password "abbhdwsy")
;(def password "abc")

(defn get-hash
  [x]
  (md5 (str password x)))

(defn get-hash-num
  [hash]
  (let [start (subs hash 0 5)]
    (when (= "00000" start) (subs hash 5 6))))

(defn search-hash
  [x pin]
  (cond
    (= 8 (count pin)) pin
    :else
      (let [next (get-hash-num (get-hash x))]
        (cond
          (nil? next) (recur (inc x) pin)
          :else  (recur (inc x) (conj pin next))))))

(defrecord Pin-pos [pos num])

(def five-zeros "00000")

(defn get-hash-num-2
  [hash]
  (when (str/starts-with? hash five-zeros)
    (let [pos (subs hash 5 6)
          num (subs hash 6 7)]
      (when (Character/isDigit (first pos))
        (->Pin-pos (Integer. (subs hash 5 6)) (subs hash 6 7))))))

(defn add-if-not-exists
  [pin new-item]
  (pprint (str "exists check: " (:pos new-item) " key: " (:num new-item)))
  (let [pos (:pos new-item)
        existing (filter #(= pos (:pos %)) pin)]
    (cond
      (> pos 7) pin
      (> (count existing) 0) pin
      :else (conj pin new-item))))

(defn search-hash-2
  [x pin]
  (when (= (mod x 1000000) 0) (println (str "solved: " (count pin) " x: " x)))
  (cond
    (= 8 (count pin)) pin
    :else
      (let [next (get-hash-num-2 (get-hash x))]
        (cond
          (nil? next) (recur (inc x) pin)
          :else  (recur (inc x) (add-if-not-exists pin next))))))

(defn solve-part-1
  [input-data]
  (search-hash 0 []))

(defn solve-part-2
  [input-data]
  (search-hash-2 0 []))
