---
title: "tiktokify: A Hackathon winning product"
description: "A post detailing my experience at the MongoDB Gen AI Hackathon and the product me and my team built over the span of 8 hours."
date: 2024-04-14
draft: false
tags: [hackathon, genAI, development, embeddings, mongoDB]
weight: 100
---


Last week, me and my team, **http418**, won the [MongoDB Gen AI Hackathon](https://lu.ma/MongoDB-GenAI-NYC). We created a tool called "tiktokify" which can automatically generate a small 30 second highlight clip from a video using clever prompt engineering and SoTA(State-of-The-Art) embedding models. We won the "Best Use of Nomic Platform" prize and here is my account of what we built and how we won.

---
### TL;DR(presented by ChatGPT): 
Our team won the MongoDB Gen AI Hackathon with "tiktokify," a tool generating 30-second highlight clips from videos. Leveraging cutting-edge embeddings and prompt engineering, we automated the creation of engaging content. We used Nomic AI's embeddings model, transcribed video audio for text summaries, and extracted video clip embeddings stored in MongoDB. Stitching together the most relevant clips, we added AI-generated narration for a complete TikTok-ready highlight. Check out our project code on GitHub and explore AWS services for similar functionality.

{{< figure src="/mongo_winning_shot.png" caption="After hours of building, the money shot" align="center">}}

---
## The problem
Going in, we knew we wanted to do something video related because let's face it, RAGs(Retrieval Augmented Generations) have already been proven to work pretty well and if you want create a high quality RAG then it becomes all about the quality of data. We wanted to try something different, something new, something video related, something which none of us had tried before. After the breakfast and some very quick talks by the sponsors of the hackathons, we started brainstorming about things we could try out in the video space.

> {{< figure src="/RAG.png" align="center" >}}
>
> <cite>Retrieval Augmented Generation(RAGs) explained[^1]</cite>

[^1]: This image was taken from Pinecone's article about RAGs. https://www.pinecone.io/learn/retrieval-augmented-generation/

We naturally looked at TikTok, Instagram Reels and YouTube Shorts and how they have blown up the short video content space. Something we realized is that the sort of stuff that people like looking at is very short 15-30 seconds summary of any and every thing, let it be programming tutorials, videos about reddit threads, football goals, cricket clips everything is TikTokable. What if we can make this content creation space automated? What if we help make the [dead internet theory](https://en.wikipedia.org/wiki/Dead_Internet_theory) come true? 

---
## The solution
We settled on creating a video summarizer, something that would take in a long video and "tiktokify" it, i.e., generate a highlight from it. To add to this, we also decided to add an AI voice narration so that the audio is not all junk and we can create a hook to our "tiktokified" video. 

---
## The approach
Before we dive into the approach, one very important thing which we used, that I feel might need some explanation, is embeddings. 

#### Embeddings
Embeddings are basically numerical representations of real-world data in a form which is easier for Neural Nets to understand. Embeddings are created in such a way that similar types of objects are placed near to each other in a latent space. Typically embeddings have been representations of only one form of data(either image, audio or text) but, we wanted a multimodal embedding model because of obvious reasons.

> {{< figure src="/Embeddings.png" align="center" >}} 
> 
> <cite>A visual aid for learning embeddings[^2]</cite>

[^2]: The image was taken from Pinecone's article about embeddings. https://www.pinecone.io/learn/vector-embeddings/

We initially were gonna use OpenAI's [CLIP](https://openai.com/research/clip) embeddings model but, [Andriy](https://www.linkedin.com/in/andriymulyar/), the CTO of [Nomic AI](https://www.nomic.ai/), told us about their embedding model which had a larger token size than CLIP by almost a scale of 100. So naturally, we decided to use [Nomic's Embedding Model](https://huggingface.co/nomic-ai/nomic-embed-text-v1.5) instead.


Okay, so now that embeddings are out of the way the basic idea of what we wanted to build is pretty simple, here are the steps: 
1. Generate some kind of text summary of the video

    a. Using some kind of video model to directly generate the summary.

    b. Extract images at specific intervals, run image captioning on them and generate a summary using the image captions.

    c. Transcribing the audio from the video and generate a summary of it.

2. Create video clip embeddings for querying
3. Use the generated text summary in the first step as a query vector and perform a cosine similarity on the video clip embeddings from step 2 
4. Get the clips and stitch them together.

#### Step 1: Generate a text summary
**Option 1** wasn't really feasible for us because we had not explored this space before and from the looks of it, this required a lot computational power which we didn't have, considering we are still students.

**Option 2** seemed really promising since we already had some experience working with [BLIP](https://arxiv.org/abs/2201.12086), an image captioning model. One thing we didn't realise was, since we are going to be feeding images from a video to BLIP, the context of the video as a whole will be lost. For example, the first few 10 frames from the video would produce the same exact output. This is far from ideal and could create problems while generating the summary, not to mention the loss of computational power which we already were short on. This option still was very promising but we didn't really have the time to explore it more.

**Option 3** fit our usecase perfectly because for our demo and testing purposes, we had chosen a football game as our input, which has a lot of commentary in it. Obviously, this would render other types of "silent" videos un-tiktokifyable but we didn't really have a choice. We realized, we could even resort to using the SRT file associated with a video instead of relying on the audio itself. We wrote a really [clever prompt](https://github.com/herzo175/mongodb-apr-2024-hackathon/blob/e9acc18128b82a824d9c22fa263695c99d7a89c6/research/create_text_summary_from_transcript.ipynb#L442) and fed the SRT file to ChatGPT and got a great summary from it.

#### Step 2: Create video clip embeddings for querying
We decided to use a very rudimentary approach to create embeddings from the video. We extracted frames, at every 0.5 second, from the video and fed them to the embedding model to generate embeddings. We stored these embeddings along with their timestamps in our MongoDB Atlas instance. 

#### Step 3 & 4: Stitching it all together
Now that we have our video embeddings and our search query vector, we just had to create an index in our MongoDB and run a cosine similarity search using our query vector. We got the top 50 most relevant images, sorted them in an ascending order and then extracted short 0.5 seconds clips starting from those timestamps. For example, let's say I have the most relevant timestamp as 12.5 second, I will extract a clip that has the range of 12.5-13 seconds, in this way we got a 25(50 frames * 0.5 sec clips) second highlight video. These clips were then merged using ffmpeg and then again, using ffmpeg, we added in an AI generated text-to-speech of the text summary in the video. 

We have a football highlight completely using AI up and running!!! 

---
## The end
Thanks for reading the story of what we built and how we won the MongoDB Gen AI Hackathon. Thanks to MongoDB for organizing this hackathon and also huge thanks for Nomic AI and Andriy for helping us out during the hackathon.

You can checkout the code for the hackathon here: https://github.com/herzo175/mongodb-apr-2024-hackathon

If you have money you can use this AWS service to achieve the same thing: https://aws.amazon.com/blogs/media/video-summarization-with-aws-artificial-intelligence-ai-and-machine-learning-ml-services/

The video we used: https://www.youtube.com/watch?v=h4m68r8kWAc

Great blog on storing and search embeddings using Nomic and MongoDB: https://blog.nomic.ai/posts/nomic-embed-mongo

---

P.S. If you're looking to hire, I have experience in building RAGs, working in systems, building all sorts of random things. You can have a look at my [resume](https://resume.akshatsharma.xyz) or contact me using [LinkedIn](https://linkedin.com/in/akshat-sharma-2602), [Email](mailto:akshatsharma2602@gmail.com) or [Twitter/X](https://x.com/akshat2602) and checkout my [GitHub](https://github.com/akshat2602)