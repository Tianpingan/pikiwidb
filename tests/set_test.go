/*
 * Copyright (c) 2023-present, Qihoo, Inc.  All rights reserved.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

package pikiwidb_test

import (
	"context"
	"log"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redis/go-redis/v9"

	"github.com/OpenAtomFoundation/pikiwidb/tests/util"
)

var _ = Describe("Set", Ordered, func() {
	var (
		ctx    = context.TODO()
		s      *util.Server
		client *redis.Client
	)

	// BeforeAll closures will run exactly once before any of the specs
	// within the Ordered container.
	BeforeAll(func() {
		config := util.GetConfPath(false, 0)

		s = util.StartServer(config, map[string]string{"port": strconv.Itoa(7777)}, true)
		Expect(s).NotTo(Equal(nil))
	})

	// AfterAll closures will run exactly once after the last spec has
	// finished running.
	AfterAll(func() {
		err := s.Close()
		if err != nil {
			log.Println("Close Server fail.", err.Error())
			return
		}
	})

	// When running each spec Ginkgo will first run the BeforeEach
	// closure and then the subject closure.Doing so ensures that
	// each spec has a pristine, correctly initialized, copy of the
	// shared variable.
	BeforeEach(func() {
		client = s.NewClient()
		Expect(client.FlushDB(ctx).Err()).NotTo(HaveOccurred())
		time.Sleep(1 * time.Second)
	})

	// nodes that run after the spec's subject(It).
	AfterEach(func() {
		err := client.Close()
		if err != nil {
			log.Println("Close client conn fail.", err.Error())
			return
		}
	})

	//TODO(dingxiaoshuai) Add more test cases.
	It("SUnion", func() {
		sAdd := client.SAdd(ctx, "set1", "a")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "b")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sAdd = client.SAdd(ctx, "set2", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set2", "d")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set2", "e")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sUnion := client.SUnion(ctx, "set1", "set2")
		Expect(sUnion.Err()).NotTo(HaveOccurred())
		Expect(sUnion.Val()).To(HaveLen(5))

		sUnion = client.SUnion(ctx, "nonexistent_set1", "nonexistent_set2")
		Expect(sUnion.Err()).NotTo(HaveOccurred())
		Expect(sUnion.Val()).To(HaveLen(0))

		//del
		del := client.Del(ctx, "set1", "set2")
		Expect(del.Err()).NotTo(HaveOccurred())
	})

	It("should SUnionStore", func() {
		sAdd := client.SAdd(ctx, "set1", "a")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "b")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sAdd = client.SAdd(ctx, "set2", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set2", "d")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set2", "e")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sUnionStore := client.SUnionStore(ctx, "set", "set1", "set2")
		Expect(sUnionStore.Err()).NotTo(HaveOccurred())
		Expect(sUnionStore.Val()).To(Equal(int64(5)))

		//sMembers := client.SMembers(ctx, "set")
		//Expect(sMembers.Err()).NotTo(HaveOccurred())
		//Expect(sMembers.Val()).To(HaveLen(5))

		//del
		del := client.Del(ctx, "set1", "set2", "set")
		Expect(del.Err()).NotTo(HaveOccurred())
	})
	It("Cmd SADD", func() {
		log.Println("Cmd SADD Begin")
		Expect(client.SAdd(ctx, "myset", "one", "two").Val()).NotTo(Equal("FooBar"))
	})
	It("should SAdd", func() {
		sAdd := client.SAdd(ctx, "setSAdd1", "Hello")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		sAdd = client.SAdd(ctx, "setSAdd1", "World")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		sAdd = client.SAdd(ctx, "setSAdd1", "World")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(0)))

		// sMembers := client.SMembers(ctx, "set")   After the smember command is developed, uncomment it to test smember command.
		// Expect(sMembers.Err()).NotTo(HaveOccurred())
		// Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))

		//del
		del := client.Del(ctx, "setSAdd1")
		Expect(del.Err()).NotTo(HaveOccurred())
	})

	It("should SAdd strings", func() {
		set := []string{"Hello", "World", "World"}
		sAdd := client.SAdd(ctx, "setSAdd2", set)
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(2)))

		// sMembers := client.SMembers(ctx, "set") After the smember command is developed, uncomment it to test smember command.
		// Expect(sMembers.Err()).NotTo(HaveOccurred())
		// Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))
		//del
		del := client.Del(ctx, "setSAdd2")
		Expect(del.Err()).NotTo(HaveOccurred())
	})
	It("should SInter", func() {
		sAdd := client.SAdd(ctx, "set1", "a")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "b")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sAdd = client.SAdd(ctx, "set2", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set2", "d")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set2", "e")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sInter := client.SInter(ctx, "set1", "set2")
		Expect(sInter.Err()).NotTo(HaveOccurred())
		Expect(sInter.Val()).To(Equal([]string{"c"}))

		sInter = client.SInter(ctx, "nonexistent_set1", "nonexistent_set2")
		Expect(sInter.Err()).NotTo(HaveOccurred())
		Expect(sInter.Val()).To(HaveLen(0))

		//del
		del := client.Del(ctx, "set1", "set2")
		Expect(del.Err()).NotTo(HaveOccurred())
	})

	It("should SInterStore", func() {
		sAdd := client.SAdd(ctx, "set1", "a")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "b")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sAdd = client.SAdd(ctx, "set2", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set2", "d")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set2", "e")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sInterStore := client.SInterStore(ctx, "set", "set1", "set2")
		Expect(sInterStore.Err()).NotTo(HaveOccurred())
		Expect(sInterStore.Val()).To(Equal(int64(1)))

		// sMembers := client.SMembers(ctx, "set")  // After the smember command is developed, uncomment it to test command.
		// Expect(sMembers.Err()).NotTo(HaveOccurred())
		// Expect(sMembers.Val()).To(Equal([]string{"c"}))
		//del
		del := client.Del(ctx, "set1", "set2", "set")
		Expect(del.Err()).NotTo(HaveOccurred())
	})

	It("should SCard", func() {
		sAdd := client.SAdd(ctx, "setScard", "Hello")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		sAdd = client.SAdd(ctx, "setScard", "World")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		sCard := client.SCard(ctx, "setScard")
		Expect(sCard.Err()).NotTo(HaveOccurred())
		Expect(sCard.Val()).To(Equal(int64(2)))
        })


        It("should SPop", func() {
                sAdd := client.SAdd(ctx, "setSpop", "one")
                Expect(sAdd.Err()).NotTo(HaveOccurred())
                sAdd = client.SAdd(ctx, "setSpop", "two")
                Expect(sAdd.Err()).NotTo(HaveOccurred())
                sAdd = client.SAdd(ctx, "setSpop", "three")
                Expect(sAdd.Err()).NotTo(HaveOccurred())
                sAdd = client.SAdd(ctx, "setSpop", "four")
                Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSpop", "five")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sPopN := client.SPopN(ctx, "setSpop", 3)
		Expect(sPopN.Err()).NotTo(HaveOccurred())
		Expect(sPopN.Val()).To(HaveLen(3))
		/*
		sMembers := client.SMembers(ctx, "setSpop")
		Expect(sMembers.Err()).NotTo(HaveOccurred())
		Expect(sMembers.Val()).To(HaveLen(2))
		*/


	})

	It("should SMove", func() {
		sAdd := client.SAdd(ctx, "set1", "one")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set1", "two")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sAdd = client.SAdd(ctx, "set2", "three")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sMove := client.SMove(ctx, "set1", "set2", "two")
		Expect(sMove.Err()).NotTo(HaveOccurred())
		Expect(sMove.Val()).To(Equal(true))

		sIsMember := client.SIsMember(ctx, "set1", "two")
		Expect(sIsMember.Err()).NotTo(HaveOccurred())
		Expect(sIsMember.Val()).To(Equal(false))

		sIsMember = client.SIsMember(ctx, "set2", "two")
		Expect(sIsMember.Err()).NotTo(HaveOccurred())
		Expect(sIsMember.Val()).To(Equal(true))
	})

	It("should SRem", func() {
		sAdd := client.SAdd(ctx, "set", "one")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set", "two")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "set", "three")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sRem := client.SRem(ctx, "set", "one")
		Expect(sRem.Err()).NotTo(HaveOccurred())
		Expect(sRem.Val()).To(Equal(int64(1)))

		sRem = client.SRem(ctx, "set", "four")
		Expect(sRem.Err()).NotTo(HaveOccurred())
		Expect(sRem.Val()).To(Equal(int64(0)))

 		// sMembers := client.SMembers(ctx, "set")
 		// Expect(sMembers.Err()).NotTo(HaveOccurred())
 		// Expect(sMembers.Val()).To(ConsistOf([]string{"three", "two"}))
	})

	It("should SRandmember", func() {
	    	sAdd := client.SAdd(ctx, "set", "one")
	    	Expect(sAdd.Err()).NotTo(HaveOccurred())
	    	sAdd = client.SAdd(ctx, "set", "two")
	    	Expect(sAdd.Err()).NotTo(HaveOccurred())
	    	sAdd = client.SAdd(ctx, "set", "three")
	    	Expect(sAdd.Err()).NotTo(HaveOccurred())

	    	member, err := client.SRandMember(ctx, "set").Result()
	    	Expect(err).NotTo(HaveOccurred())
	    	Expect(member).NotTo(Equal(""))

	    	members, err := client.SRandMemberN(ctx, "set", 2).Result()
	    	Expect(err).NotTo(HaveOccurred())
		Expect(members).To(HaveLen(2))
    	})

	It("should SMembers", func() {
		sAdd := client.SAdd(ctx, "setSMembers", "Hello")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSMembers", "World")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sMembers := client.SMembers(ctx, "setSMembers")
		Expect(sMembers.Err()).NotTo(HaveOccurred())
		Expect(sMembers.Val()).To(ConsistOf([]string{"Hello", "World"}))
	})

	It("should SDiff", func() {
		sAdd := client.SAdd(ctx, "setSDiff1", "a")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSDiff1", "b")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSDiff1", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sAdd = client.SAdd(ctx, "setSDiff2", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSDiff2", "d")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSDiff2", "e")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sDiff := client.SDiff(ctx, "setSDiff1", "setSDiff2")
		Expect(sDiff.Err()).NotTo(HaveOccurred())
		Expect(sDiff.Val()).To(ConsistOf([]string{"a", "b"}))

		sDiff = client.SDiff(ctx, "nonexistent_setSDiff1", "nonexistent_setSDiff2")
		Expect(sDiff.Err()).NotTo(HaveOccurred())
		Expect(sDiff.Val()).To(HaveLen(0))
	})

	It("should SDiffstore", func() {
		sAdd := client.SAdd(ctx, "setSDiffstore1", "a")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSDiffstore1", "b")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSDiffstore1", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sAdd = client.SAdd(ctx, "setSDiffstore2", "c")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSDiffstore2", "d")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		sAdd = client.SAdd(ctx, "setSDiffstore2", "e")
		Expect(sAdd.Err()).NotTo(HaveOccurred())

		sDiffStore := client.SDiffStore(ctx, "setKey", "setSDiffstore1", "setSDiffstore2")
		Expect(sDiffStore.Err()).NotTo(HaveOccurred())
		Expect(sDiffStore.Val()).To(Equal(int64(2)))

		sMembers := client.SMembers(ctx, "setKey")
		Expect(sMembers.Err()).NotTo(HaveOccurred())
		Expect(sMembers.Val()).To(ConsistOf([]string{"a", "b"}))
	})

	It("should SScan", func() {
		// add elements first
		sAdd := client.SAdd(ctx, "setSScan1", "user1")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		sAdd = client.SAdd(ctx, "setSScan1", "user2")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		sAdd = client.SAdd(ctx, "setSScan1", "user3")
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(1)))

		set := []string{"Hello", "World", "World"}
		sAdd = client.SAdd(ctx, "setSScan1", set)
		Expect(sAdd.Err()).NotTo(HaveOccurred())
		Expect(sAdd.Val()).To(Equal(int64(2)))

		// func (c Client) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
		sScan:=client.SScan(ctx,"setSScan1",0,"*",5)
		Expect(sScan.Err()).NotTo(HaveOccurred())
		Expect(sScan.Val()).To(ConsistOf([]string{"user1", "user2","user3","Hello","World"}))
		
		sScan=client.SScan(ctx,"setSScan1",0,"user*",5)
		Expect(sScan.Err()).NotTo(HaveOccurred())
		Expect(sScan.Val()).To(ConsistOf([]string{"user1", "user2","user3"}))

		sScan=client.SScan(ctx,"setSScan1",0,"He*",5)
		Expect(sScan.Err()).NotTo(HaveOccurred())
		Expect(sScan.Val()).To(ConsistOf([]string{"Hello"}))
		
		// sScan=client.SScan(ctx,"setSScan1",0,"*",-1)
		// Expect(sScan.Err()).To(HaveOccurred())
		

		//del
		del := client.Del(ctx, "setSScan1")
		Expect(del.Err()).NotTo(HaveOccurred())
	})



})
