import type {NextPage} from 'next';
import {useEffect, useRef, useState} from 'react';
import Head from 'next/head';
import {Button, Center, Container, Heading, Input, Spinner, Text} from '@chakra-ui/react';
import About from '../components/About';
import Instruction from '../components/Instruction';
import axios from "axios";

type MessageStats = {
    totalCount: number;
    list: {
        name: string;
        count: number;
    }[];
}

const Home: NextPage = () => {
    const [token, setToken] = useState<string>('');
    const [chatId, setChatId] = useState<string>('');

    const [stats, setStats] = useState<MessageStats>(null);
    const [loading, setLoading] = useState<boolean>(false);

    const statsEl = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const defaultToken = window.localStorage.getItem('token');
        if (defaultToken) {
            setToken(defaultToken);
        }

        const defaultChatId = window.localStorage.getItem('chatId');
        if (defaultChatId) {
            setChatId(defaultChatId);
        }
    }, [])

    const onChangeToken = (e) => {
        setToken(e.target.value)
        window.localStorage.setItem('token', e.target.value);
    }

    const onChangeChatId = (e) => {
        setChatId(e.target.value)
        window.localStorage.setItem('chatId', e.target.value);
    }

    const getStats = async () => {
        if (loading) {
            return;
        }

        if (stats) {
            setStats(null);
        }

        try {
            if (!token || !chatId) {
                alert('Empty token or chat id');
                return;
            }

            const chatIdInt = parseInt(chatId);
            if (isNaN(chatIdInt)) {
                alert('invalid chat id');
                return;
            }

            setLoading(true);
            const response = await axios.get('http://localhost:8080/stats', {
                params: {
                    token,
                    chat_id: chatIdInt,
                }
            });

            setStats(response.data);
        } catch (e) {
            alert(`${e.response.status} error code\nError text: ${e.response.data.message}`);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (statsEl.current) {
            statsEl.current.scrollIntoView({ block: "start", behavior: 'smooth' });
        }
    }, [stats, statsEl]);

    return (
        <div>
            <Head>
                <title>VK Message Counter</title>
                <meta name="description" content="Simple app for getting vk message statistics"/>
                <link rel="icon" href="/favicon.ico"/>
            </Head>

            <Container mt="4" maxW="container.sm">
                <About/>
                <Instruction/>

                <Heading mb="4">Do it!</Heading>
                <Input placeholder="API Token" mb="4" defaultValue={token} onChange={onChangeToken}/>
                <Input placeholder="Chat ID" mb="4" defaultValue={chatId} onChange={onChangeChatId}/>
                <Button w="full" mb="4" onClick={getStats}>Go</Button>

                {loading && <>
                    <Center mt="4">
                        <Spinner/>
                    </Center>
                </>}

                {stats &&
                    <div ref={statsEl}>
                        <Heading mb="4">Results</Heading>
                        <Text lineHeight="7" mb="4">
                            <b>Total count - {stats.totalCount}</b>
                            <br/>
                            {stats.list.map((item) => (
                                <>
                                    {item.name} - {item.count}
                                    <br/>
                                </>
                            ))}
                        </Text>
                    </div>
                }
            </Container>
        </div>
    );
};

export default Home;
