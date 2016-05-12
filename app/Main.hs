{-# LANGUAGE OverloadedStrings #-}
import Data.Char (isPunctuation, isSpace)
import Data.Monoid (mappend)
import Data.Text (Text)
import Control.Exception (finally)
import Control.Monad (forM_, forever)
import Control.Concurrent (MVar, newMVar, modifyMVar_, modifyMVar, readMVar)
import qualified Data.Text as T
import qualified Data.Text.IO as T

import qualified Network.WebSockets as WS

type Client = WS.Connection

type ServerState = [Client]

newServerState :: ServerState
newServerState = []

addClient :: Client -> ServerState -> ServerState
addClient client clients = client : clients

removeClient :: Client -> ServerState -> ServerState
removeClient client = id -- filter (/= client)

broadcast :: Text -> ServerState -> IO ()
broadcast message clients = do
  T.putStrLn message
  forM_ clients $ \conn -> WS.sendTextData conn message

main :: IO ()
main = do
  state <- newMVar newServerState
  WS.runServer "0.0.0.0" 8080 $ application state

application :: MVar ServerState -> WS.ServerApp
application state pending = do
  conn <- WS.acceptRequest pending
  WS.forkPingThread conn 30
  msg <- WS.receiveData conn
  clients <- readMVar state
  case T.take 8 msg of
    "join    " -> flip finally disconnect $ do
      modifyMVar_ state $ \s -> do
        let s' = addClient client s
        broadcast "someone joined" s'
        return s'
      talk conn state
      where
        client     = conn
        disconnect = do
          s <- modifyMVar state $ \s -> let s' = removeClient client s in return (s', s')
          broadcast "someone disconnected" s
    otherwise  -> do WS.sendTextData conn $ "Received message '" `mappend` msg `mappend` "'"

talk :: WS.Connection -> MVar ServerState -> IO ()
talk conn state = forever $ do
  msg <- WS.receiveData conn
  case T.take 8 msg of
    "update  " -> undefined
    "end turn" -> undefined
    otherwise  -> readMVar state >>= broadcast ("Received message '" `mappend` msg `mappend` "'")
