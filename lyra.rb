class Lyra < Formula
  desc "a lightweight encryption tool that makes protecting your sensitive files easy"
  homepage "https://github.com/azohra/lyra"
  url "https://github.com/azohra/lyra/releases/download/v1.0.0/lyra_darwin_1.0.0.tar.gz"
  sha256 "7978aa39a38683de94770eb38566c68bcbfe166b94f166d86cc39e4a82d6e7f2"

  bottle :unneeded

  def install
    bin.install "lyra"
  end

  test do
    system "#{bin}/lyra", "-h"
  end
end
